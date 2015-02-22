package obj

type Material [3]uint

type MaterialType int

const (
	MATERIAL_AMBIENT MaterialType = iota
	MATERIAL_DIFFUSE
	MATERIAL_SPECULAR
)

// I need to specify which wavelength of light are reflected and which are held back by the material.
type material struct {
	Reflection uint8
	Ambient    *Material
	Diffuse    *Material
	Specular   *Material
}
