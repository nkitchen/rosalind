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

// ReadFasta reads a single String and assumes that nothing else follows it
// in the reader's data.
func ReadFasta(r io.Reader) (ret String, retErr error) {
	lr := lineReader{bufio.NewReader(r), 0}

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
		line, err := lr.ReadLine()

		line = strings.TrimRight(line, "\n")
		data = append(data, line...)

		if err == io.EOF || line == "" || line[0] == '>' {
			return String{Description: desc, Data: string(data)}, nil
		}
		if err != nil {
			return String{Description: desc, Data: string(data)},
			       newFastaError(err.Error())
		}
	}
	return
}
