package httpsignatures

import (
	"net/http"
	"testing"
)

const testErrDigestType = "*httpsignatures.ErrDigest"

func TestVerifyDigest(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Valid MD5 digest",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "MD5=Sd/dVLAcvNLSq16eXua5uQ=="),
			},
			want:        true,
			wantErrType: testErrDigestType,
			wantErrMsg:  "",
		},
		{
			name: "Valid SHA-256 digest",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "SHA-256=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE="),
			},
			want:        true,
			wantErrType: testErrDigestType,
		},
		{
			name: "Valid SHA-512 digest",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "SHA-512=WZDPaVn/7XgHaAy8pmojAkGWoRx2UFChF41A2svX+TaPm+"+
					"AbwAgBWnrIiYllu7BNNyealdVLvRwEmTHWXvJwew=="),
			},
			want:        true,
			wantErrType: testErrDigestType,
		},
		{
			name: "Invalid MD5 digest (decode error)",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "MD5=123456"),
			},
			want:        false,
			wantErrType: testErrDigestType,
			wantErrMsg:  "ErrDigest: error decode digest from base64: illegal base64 data at input byte 4",
		},
		{
			name: "Invalid MD5 wrong digest",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "MD5=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE="),
			},
			want:        false,
			wantErrType: testErrDigestType,
			wantErrMsg:  "ErrDigest: wrong digest: ErrCrypto: wrong hash",
		},
		{
			name: "Invalid digest header",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "SHA-512="),
			},
			want:        false,
			wantErrType: testErrParserType,
			wantErrMsg:  "ErrParser: empty digest value",
		},
		{
			name: "Unsupported digest hash algorithm",
			args: args{
				r: testGetDigestRequestFunc(testBodyExample, "SHA-0=test"),
			},
			want:        false,
			wantErrType: testErrDigestType,
			wantErrMsg:  "ErrDigest: unsupported digest hash algorithm 'SHA-0'",
		},
		{
			name: "Empty body",
			args: args{
				r: testGetDigestRequestFunc("", "MD5=xxx"),
			},
			want:        false,
			wantErrType: testErrDigestType,
			wantErrMsg:  "ErrDigest: empty body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDigest()
			err := d.Verify(tt.args.r)
			got := err == nil
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestCreateDigest(t *testing.T) {
	type args struct {
		alg string
		r   *http.Request
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Valid MD5 digest",
			args: args{
				alg: "MD5",
				r:   testGetDigestRequestFunc(testBodyExample, ""),
			},
			want:        "MD5=Sd/dVLAcvNLSq16eXua5uQ==",
			wantErrType: "",
			wantErrMsg:  "",
		},
		{
			name: "Unsupported digest alg",
			args: args{
				alg: "MD4",
				r:   testGetDigestRequestFunc(testBodyExample, ""),
			},
			want:        "",
			wantErrType: testErrDigestType,
			wantErrMsg:  "ErrDigest: unsupported digest hash algorithm 'MD4'",
		},
		{
			name: "Create digest error",
			args: args{
				alg: "ERR",
				r:   testGetDigestRequestFunc(testBodyExample, ""),
			},
			want:        "",
			wantErrType: testErrDigestType,
			wantErrMsg:  "ErrDigest: error creating digest hash 'ERR': create hash error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDigest()
			d.SetDigestHashAlgorithm(testErrAlg{})
			got, err := d.Create(tt.args.alg, tt.args.r)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestDigestSetDigestHashAlgorithm(t *testing.T) {
	tests := []struct {
		name string
		arg  DigestHashAlgorithm
	}{
		{
			name: "Set new algorithm OK",
			arg:  testAlg{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDigest()
			d.SetDigestHashAlgorithm(tt.arg)
			if _, ok := d.alg[testAlgName]; ok == false {
				t.Error("algorithm not found")
			}
		})
	}
}

func TestDigestSetDigestDefaultHashAlgorithm(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "Set new default algorithm OK",
			arg:  algSha256,
			want: true,
		},
		{
			name: "Algorithm was not set",
			arg:  "test",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDigest()
			_ = d.SetDefaultDigestHashAlgorithm(tt.arg)
			got := false
			if d.defaultAlg == tt.arg {
				got = true
			}
			if got != tt.want {
				t.Error("default algorithm was not set")
			}
		})
	}
}
