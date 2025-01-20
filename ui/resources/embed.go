package resources

import (
	_ "embed"
)

var (
	//go:embed Nw.png
	Nw_png []byte

	//go:embed Bw.png
	Bw_png []byte

	//go:embed Rw.png
	Rw_png []byte

	//go:embed Qw.png
	Qw_png []byte

	//go:embed Kw.png
	Kw_png []byte

	//go:embed Pw.png
	Pw_png []byte

	//go:embed Nb.png
	Nb_png []byte

	//go:embed Bb.png
	Bb_png []byte

	//go:embed Rb.png
	Rb_png []byte

	//go:embed Qb.png
	Qb_png []byte

	//go:embed Kb.png
	Kb_png []byte

	//go:embed Pb.png
	Pb_png []byte
)

var PiecePngs [12]*[]byte

func init() {
	PiecePngs = [12]*[]byte{&Nw_png, &Bw_png, &Rw_png, &Qw_png, &Kw_png, &Pw_png, &Nb_png, &Bb_png, &Rb_png, &Qb_png, &Kb_png, &Pb_png}
}
