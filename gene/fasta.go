package gene

import "bufio"
import "fmt"
import "io"
import "strings"

type String struct {
	Description string
	Data string
}

type FastaError string
func (f FastaError) Error() string {
	return string(f)
}

type lineReader struct {
	*bufio.Reader
	Lineno int
}

func (r *lineReader) ReadLine() (string, error) {
	line, err := r.ReadString('\n')
	r.Lineno++
	return line, err
}

func (lr *lineReader) ReadFasta() (ret String, retErr error) {
	newFastaError := func(s string) FastaError {
		return FastaError(fmt.Sprintf("Line %v: %s", lr.Lineno, s))
	}

	line, err := lr.ReadLine()
	if err == io.EOF {
		retErr = FastaError("EOF before description")
		return
	}
	if err != nil {
		retErr = newFastaError(err.Error())
		return
	}

	line = strings.TrimRight(line, "\n")

	if line[0] != '>' {
		retErr = newFastaError("Missing > before description")
		return
	}

	desc := line[1:]
	data := []byte{}
	for {
		buf, err := lr.Reader.Peek(1)
		done := len(buf) == 0 && err == io.EOF
		if len(buf) > 0 && (buf[0] == '\n' || buf[0] == '>') {
			done = true
		}
		if done {
			return String{Description: desc, Data: string(data)}, err
		}

		line, err := lr.ReadLine()
		if err != nil {
			return String{Description: desc, Data: string(data)},
			       newFastaError(err.Error())
		}

		line = strings.TrimRight(line, "\n")
		data = append(data, line...)
	}
	return
}

// ReadFasta reads a single String and assumes that nothing else follows it
// in the reader's data.
func ReadFasta(r io.Reader) (String, error) {
	lr := lineReader{bufio.NewReader(r), 0}
	return lr.ReadFasta()
}

func ReadAllFasta(r io.Reader) ([]String, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	lr := lineReader{br, 0}
	a := []String{}
	for {
		s, err := lr.ReadFasta()
		a = append(a, s)
		if err == io.EOF {
			return a, nil
		}
		if err != nil {
			return a, err
		}
	}
	return nil, nil
}

func WriteFasta(w io.Writer, s String) error {
	_, err := fmt.Fprintf(w, ">%s\n", s.Description)
	if err != nil {
		return err
	}
	for i := 0; i < len(s.Data); i += 80 {
		t := s.Data[i:]
		if len(t) > 80 {
			t = t[:80]
		}
		_, err = fmt.Fprintln(w, t)
		if err != nil {
			return err
		}
	}
	return nil
}
