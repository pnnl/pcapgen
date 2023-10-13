#! /bin/sh
# Cyber Fire Foundry: Netarch 25000 solution
# 2020 Neale Pickett <neale@lanl.gov>
# Public Domain
#
# This is not an efficient solution.
#
# SPOILER ALERT
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
#
# DUMBLEDORE DIES


die () {
        echo "$*" 1>&2
        exit 1
}

assert () {
        test "$@" || die "Assert: $*"
}

hex2dec () {
	printf "%d" "0x$1"
}

mkdir -p xfer
frame=0
pmerge "$@" | pcat | while read ts proto opts src dst payload; do
        frame=$(($frame + 1))
        when=$(TZ=Z date -d @${ts%.*} "+%Y-%m-%d %H:%M:%S")
        opcode=None
        opcode_desc=None

        # ICMP sequence number is not always 0
        # assert $(echo $payload | slice 0 8) = 00000000
        payload=$(echo $payload | slice 8 | unhex | xor -x 70 65 67 6d 0a 53 45 5f  0a 4d 45 5e 0a 43 5e 0b | hex | tr -d ' ')

        if [ $frame -le 3 ]; then
		opcode="Handshake"
        elif [ -z "$payload" ]; then
		opcode="Empty"
        elif [ "$payload" = 00 ]; then
		opcode="ACK"
	else
		opcode=$(echo $payload | slice 0 2)
		assert $(echo $payload | slice 2 4) = 00
		session=$(echo $payload | slice 4 6)
		assert $(echo $payload | slice 6 8) = 00
		payload=$(echo $payload | slice 8)
	fi
        case $opcode in
                01)
                        opcode_desc="Filename"
                        wat=$(echo $payload | slice 0 8)
                        len=$(hex2dec $(echo $payload | slice 8 10))
                        payload=$(echo $payload | slice 10)
                        filename=$(echo $payload | unhex)
                        assert $len = $(echo $payload | unhex | wc -c)
                        : > xfer/$session
                        ln -f xfer/$session xfer/$filename
                        ;;
                02)
                        opcode_desc="Xfer"
                        len=$(hex2dec $(echo $payload | slice 2 4)$(echo $payload | slice 0 2))
                        payload=$(echo $payload | slice 4)
                        assert $len = $(echo $payload | unhex | wc -c)
                        echo $payload | unhex >> xfer/$session
                        ;;
        esac


        printf "Packet %d %s %s: %s\n" $frame "$proto" "$opcode" "$opcode_desc"
        printf "    %s -> %s (%s)\n" ${src%,*} ${dst%,*} "$when"
        printf "    Session: %s\n" "$session"
        printf "    Wat: %s\n" "$wat"
        printf "    Len: %s\n" "$len"
        printf "    Payload Length: %02x\n" $(echo -n $payload | unhex | wc -c)

        echo $payload | unhex | hd
        echo
done
