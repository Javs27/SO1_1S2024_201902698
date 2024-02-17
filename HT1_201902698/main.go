package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

const filename = "cursos.bin"

type Course struct {
	Type bool
	ID   int64
	Code int64
	Name string
}

func createFile() error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func getLastCourseID() int64 {
	file, err := os.Open(filename)
	if err != nil {
		return 0 // Si hay un error (por ejemplo, el archivo no existe), comenzamos desde 0
	}
	defer file.Close()

	// Ir al final del archivo menos el tamaño de un registro
	_, err = file.Seek(-33, io.SeekEnd)
	if err != nil {
		return 0 // Si hay un error (archivo demasiado corto), comenzamos desde 0
	}

	var lastID int64
	err = binary.Read(file, binary.LittleEndian, &lastID)
	if err != nil {
		return 0
	}

	return lastID
}

func writeCourseToFile(course Course) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Preparar un writer bufferizado para el archivo
	writer := bufio.NewWriter(file)

	// Escribir ID y Code como Varint
	err = binary.Write(writer, binary.LittleEndian, int64(course.ID))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, int64(course.Code))
	if err != nil {
		return err
	}

	// Asegurar que el nombre tenga exactamente 16 bytes
	nameBuffer := []byte(course.Name)
	if len(nameBuffer) > 16 {
		nameBuffer = nameBuffer[:16] // Cortar si es necesario
	} else {
		// Rellenar con ceros si es más corto
		for len(nameBuffer) < 16 {
			nameBuffer = append(nameBuffer, 0)
		}
	}

	// Escribir el nombre
	_, err = writer.Write(nameBuffer)
	if err != nil {
		return err
	}

	// Escribir el tipo
	typeByte := byte(0)
	if course.Type {
		typeByte = 1
	}
	_, err = writer.Write([]byte{typeByte})
	if err != nil {
		return err
	}

	// Asegurarse de que todo se ha escrito al archivo
	return writer.Flush()
}

func registerCourse() {
	var course Course

	// Obtener el último ID y autoincrementar
	course.ID = getLastCourseID() + 1

	fmt.Print("Ingrese el tipo de registro (obligatorio/no obligatorio): ")
	var tipo string
	fmt.Scan(&tipo)
	course.Type = strings.ToLower(tipo) == "obligatorio"

	fmt.Print("Ingrese el código: ")
	fmt.Scan(&course.Code)
	fmt.Print("Ingrese el nombre (máximo 16 caracteres): ")
	fmt.Scan(&course.Name)

	err := writeCourseToFile(course)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	} else {
		fmt.Println("Curso registrado con éxito.")
		fmt.Println("")
	}
}

func viewRecords() {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	var course Course
	for {
		// Leer ID como int64
		err = binary.Read(file, binary.LittleEndian, &course.ID)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error al leer el ID:", err)
			}
			break // Fin del archivo o error
		}

		// Leer Code como int64
		err = binary.Read(file, binary.LittleEndian, &course.Code)
		if err != nil {
			fmt.Println("Error al leer el Code:", err)
			break // Error
		}

		// Leer Name (16 bytes fijos)
		nameBuffer := make([]byte, 16)
		_, err = file.Read(nameBuffer)
		if err != nil {
			fmt.Println("Error al leer el Name:", err)
			break // Error
		}
		course.Name = strings.TrimRight(string(nameBuffer), "\x00")

		// Leer Type como un byte
		typeBuffer := make([]byte, 1)
		_, err = file.Read(typeBuffer)
		if err != nil {
			fmt.Println("Error al leer el Type:", err)
			break // Error
		}
		course.Type = typeBuffer[0] == 1

		fmt.Printf("ID: %d, Code: %d, Name: %s, Type: %v\n", course.ID, course.Code, course.Name, course.Type)
	}
}

func main() {
	err := createFile()
	if err != nil {
		fmt.Println("El archivo ya existe.")
		fmt.Println("")
	}

	var choice int
	for {
		fmt.Println("Menú Principal:")
		fmt.Println("1. Registro de Curso")
		fmt.Println("2. Ver Registros")
		fmt.Println("0. Salir")
		fmt.Print("Ingrese su elección: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			registerCourse()
			fmt.Println() // Espacio adicional para mejorar la visualización
			fmt.Println("")
		case 2:
			viewRecords()
			fmt.Println() // Espacio adicional para mejorar la visualización
			fmt.Println("")
		case 0:
			fmt.Println("Saliendo...")
			fmt.Println("")
			return
		default:
			fmt.Println("Opción no válida. Inténtelo de nuevo.")
			fmt.Println("")
		}
	}
}
