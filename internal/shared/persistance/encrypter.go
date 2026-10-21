package persistance

// Encrypter xifra i desxifra secrets. Exposa dos nivells perquè la codificació
// és cosa de qui crida, no del xifrat:
//
//   - Seal/Open treballen amb bytes crus. Els fa servir qui ha de triar com
//     transportar-los (p. ex. base64 URL-safe per a un state que viatja per
//     query string, on el '+' del base64 estàndard es llegiria com un espai).
//   - Encrypt/Decrypt són la comoditat per a at-rest: retornen base64 estàndard,
//     llest per desar en una columna TEXT.
type Encrypter interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)

	Seal(plaintext []byte) ([]byte, error)
	Open(sealed []byte) ([]byte, error)
}
