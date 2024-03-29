// Package mtab parses /proc/self/mounts entries on a Linux system.
package mtab // import "github.com/artyom/mtab"

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Entry corresponds to mntent struct. See getmntent(3) manpage for further
// details.
type Entry struct {
	Fsname string // name of mounted file system
	Dir    string // file system path prefix
	Type   string // mount type
	Opts   string // mount options
	Freq   int    // dump frequency in days
	Passno int    // pass number on parallel fsck
}

const mtabFmt = `%s %s %s %s %d %d`

var escaped = strings.NewReplacer(
	`\040`, " ",
	`\011`, "\t",
	`\012`, "\n",
	`\134`, `\`,
)

// unescapeFields unescapes characters on string fields. See getmntent(3)
// manpage for details.
func unescapeFields(m *Entry) {
	for _, f := range [...]*string{&m.Fsname, &m.Dir, &m.Type, &m.Opts} {
		*f = escaped.Replace(*f)
	}
}

// Entries reads mtab entries from a given file. Usually you should use
// `/etc/fstab` or `/proc/self/mounts` as a file name.
func Entries(fname string) ([]Entry, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []Entry
READLOOP:
	for {
		e := Entry{}
		n, err := fmt.Fscanf(f, mtabFmt, &e.Fsname, &e.Dir, &e.Type, &e.Opts, &e.Freq, &e.Passno)
		switch err {
		case io.EOF:
			break READLOOP
		case nil:
		default:
			return out, err
		}
		if n != 6 {
			return out, fmt.Errorf("wrong line format (invalid number of fields)")
		}
		unescapeFields(&e)
		out = append(out, e)
	}
	return out, nil
}
