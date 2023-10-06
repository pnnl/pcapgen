package conversation

// func TestRawConversation(t *testing.T) {
// 	outf, err := os.CreateTemp("", "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer outf.Close()
// 	defer os.Remove(outf.Name())

// 	wg := sync.WaitGroup{}

// 	msg := "bloop"
// 	conv := NewRawConversation(outf)
// 	wg.Add(2)
// 	go func() {
// 		defer wg.Done()
// 		if err := conv.A.Write([]byte(msg)); err != nil {
// 			t.Fatal(err)
// 		}
// 	}()
// 	go func() {
// 		defer wg.Done()
// 		buf := make([]byte, 20)
// 		if err, length := conv.B.Read(buf); err != nil {
// 			t.Error(err)
// 		} else if length != len(msg) {
// 			t.Error("Wrong length:", length)
// 		} else if string(msg[:length]) != msg {
// 			t.Error("Wrong read bytes")
// 		}
// 	}()
// 	wg.Wait()
// }
