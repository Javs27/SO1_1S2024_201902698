package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"

	_ "github.com/go-sql-driver/mysql"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	// Crear un enrutador
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{})

	//configuracion del puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.HandleFunc("/", bienvenida).Methods("GET")
	router.HandleFunc("/ram", socketMemory)
	router.HandleFunc("/cpu", socketCpu)
	router.HandleFunc("/kill", killProcess).Methods("POST")
	router.HandleFunc("/loadCpu", loadCpu).Methods("GET")

	router.HandleFunc("/datacpu", DataCpu).Methods("GET")
	router.HandleFunc("/dataram", DataRam).Methods("GET")

	router.HandleFunc("/dataramhistorial", DataRamHistorial).Methods("GET")
	router.HandleFunc("/datacpuhistorial", DataCpuHistorial).Methods("GET")

	fmt.Println("Server started on: " + port + " port")
	http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router))

}

// prueba de funcionamiento
func bienvenida(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Bienvenido al servidor del proyecto 1 de Sistemas Operativos 1"))
}

// Retrona el json del cpu
func DataCpu(w http.ResponseWriter, r *http.Request) {
	data := getCPU()

	fmt.Print(data.Usage)
	// Inserta datos en la tabla RAM_HISTORIAL
	// esto es para la base de datos
	dataSourceName := "root:JAVS27@tcp(localhost:3306)/Proyecto1SO"
	// Abre una conexión a la base de datos
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print("Conexion exitosa a MYSQL :D")
	//defer db.Close()
	err = insertarDatos2(db, int(data.Usage))
	if err != nil {
		panic(err.Error())
	}

	// Configura la cabecera de la respuesta para indicar que se envía JSON
	w.Header().Set("Content-Type", "application/json")

	// Codifica los datos en formato JSON y los envia como respuesta
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error al enviar los datos del CPU.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Datos del CPU enviados correctamente.")
}

// Retrona el json de la ram aca esta la informacion de la ram

func DataRam(w http.ResponseWriter, r *http.Request) {

	data := getRAM()
	fmt.Println(data.Used_memory) //esto nos muesta el porcentaje de la memoria ram

	// enviamos el valor a la base de datos
	// Inserta datos en la tabla RAM_HISTORIAL
	// esto es para la base de datos
	dataSourceName := "root:JAVS27@tcp(localhost:3306)/Proyecto1SO"
	// Abre una conexión a la base de datos
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print("Conexion exitosa a MYSQL :D")
	defer db.Close()
	err = insertarDatos(db, data.Used_memory)
	if err != nil {
		panic(err.Error())
	}

	// Configura la cabecera de la respuesta para indicar que se envía JSON
	w.Header().Set("Content-Type", "application/json")

	// Codifica los datos en formato JSON y los envia como respuesta
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error al enviar los datos de la ram")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Datos de la RAM enviados correctamente.")
}

func loadCpu(response http.ResponseWriter, request *http.Request) {
	numero := 123
	for i := 0; i < 100; i++ {
		numero = numero + numero
	}
	response.Write([]byte("exito"))
}

func readerCPU(connection *websocket.Conn) {
	for {
		messageType, p, err := connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		if err := connection.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func readerRam(connection *websocket.Conn) {
	for {
		messageType, p, err := connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))

		if err := connection.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func writerRam(connection *websocket.Conn) {
	for {
		data := getRAM()
		if err := connection.WriteJSON(data); err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
func writerCpu(connection *websocket.Conn) {
	for {
		data := getCPU()
		//log.Println(data)
		if err := connection.WriteJSON(data); err != nil {

			log.Println(err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}

}

func socketMemory(response http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = func(request *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Cliente conectado: RAM")
	writerRam(ws)
	log.Println("Cliente desconectado: RAM")
}

func socketCpu(response http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = func(request *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println("BANDERA DE ERROR 2 ")
		log.Println(err)
	}
	log.Println("Cliente conectado: CPU")
	writerCpu(ws)
	log.Println("Cliente desconectado:  CPU")
}

func getCpuUsage() float64 {
	cmd := exec.Command("sh", "-c", `ps -eo pcpu | sort -k 1 -r | head -50`)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getCpuUse: ", err)
	}
	salidaAuxiliar := strings.Split(string(stdout), "\n")
	var total float64 = 0
	for i := 0; i < len(salidaAuxiliar); i++ {
		float1, _ := strconv.ParseFloat(salidaAuxiliar[i], 64)
		total += float1
	}
	total = (total / float64(len(salidaAuxiliar)-43))
	return (total)
}
func obtenerPorcentajeCPU() (float64, error) {
	// Intervalo de tiempo para calcular el uso de la CPU
	interval := 1 * time.Second

	// Obtiene las estadísticas de la CPU
	cpuPercent, err := cpu.Percent(interval, false)
	if err != nil {
		return 0.0, err
	}

	// Devuelve el porcentaje de uso de la CPU
	return cpuPercent[0], nil
}

func getCache() float64 {
	cmd := exec.Command("sh", "-c", `free -m | head -n2 | tail -1 | awk '{print $6}'`)
	stdout, err := cmd.Output()
	if err != nil {
		//fmt.Println("Error getCache: ", err)
	}
	salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
	valor, _ := strconv.ParseFloat(salida, 64)
	return valor
}

func getRAM() Memoria {
	ram, _ := ioutil.ReadFile("/proc/ram_so1_1s2024")
	var memoria Memoria
	json.Unmarshal(ram, &memoria)
	//log.Println(memoria)
	memoria.Cache_memory = getCache()
	//memoria.Used_memory = (memoria.Total_memory - memoria.Free_memory - int(getCache())) * 100 / memoria.Total_memory
	memoria.Available_memory = memoria.Free_memory + int(getCache())
	memoria.MB_memory = (memoria.Total_memory - memoria.Free_memory - int(getCache()))
	return memoria
}

func getCPU() CpuSend {
	processes, _ := ioutil.ReadFile("/proc/modulo_cpu")
	var cpu Cpu
	var cpuSend CpuSend
	json.Unmarshal(processes, &cpu)
	//fmt.Println(cpu.Usage)
	cpu.Usage, _ = obtenerPorcentajeCPU()
	//cpu.Usage = getCpuUsage()
	//fmt.Println(cpu.Usage)
	hashmap := make(map[int]string)
	var keys []int
	for i := 0; i < len(cpu.Processes); i++ {
		inputProcess := cpu.Processes[i]
		if !contains(keys, inputProcess.User) {
			keys = append(keys, inputProcess.User)
			hashmap[inputProcess.User] = getUser2(inputProcess.User)
		}
		//auxiliar := ProcessSend{Pid: inputProcess.Pid, Name: inputProcess.Name, User: strconv.Itoa(inputProcess.User), State: inputProcess.State, Ram: inputProcess.Ram, Child: inputProcess.Child}
		auxiliar := ProcessSend{Pid: inputProcess.Pid, Name: inputProcess.Name, User: hashmap[inputProcess.User], State: inputProcess.State, Ram: inputProcess.Ram, Child: inputProcess.Child}
		cpuSend.Processes = append(cpuSend.Processes, auxiliar)
	}
	cpuSend.Running = cpu.Running
	cpuSend.Sleeping = cpu.Sleeping
	cpuSend.Stopped = cpu.Stopped
	cpuSend.Total = cpu.Total
	cpuSend.Usage = cpu.Usage
	cpuSend.Zombie = cpu.Zombie
	return cpuSend
}

func getUser(nombre int) string {
	cmd := exec.Command("sh", "-c", `id -nu `+strconv.Itoa(nombre))
	stdout, err := cmd.Output()
	if err != nil {
		//fmt.Println("Error getUser: ", err)
	}
	salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
	return salida
}
func getUser2(uid int) string {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return strconv.Itoa(uid)
	}
	return u.Username
}

type JsonRequest struct {
	Pid int `json:"pid"`
}

func killProcess(response http.ResponseWriter, request *http.Request) {
	data, errRead := ioutil.ReadAll(request.Body)
	if errRead != nil {
		fmt.Println("Error al leer la solicitud: ", errRead)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"value": false, "error": "Error al leer la solicitud"}`))
		return
	}
	var requestJson JsonRequest
	err := json.Unmarshal(data, &requestJson)
	if err != nil {
		fmt.Println("Error al analizar el JSON: ", err)
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(`{"value": false, "error": "Error al analizar el JSON"}`))
		return
	}
	fmt.Println("PID a matar: ", requestJson.Pid)
	cmd := exec.Command("sh", "-c", "kill "+strconv.Itoa(requestJson.Pid))
	resultdo, err := cmd.CombinedOutput()
	fmt.Println("Resultado del comando :", string(resultdo))
	if err != nil {
		fmt.Println("Error al matar el proceso:", err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"value": false, "error": "Error al matar el proceso"}`))
		return
	}

	fmt.Println("Proceso eliminado exitosamente")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(`{"value": true}`))
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//Structs

type Memoria struct {
	Total_memory     int     `json: "total_memory"`
	Free_memory      int     `json: "free_memory"`
	Used_memory      int     `json: "used_memory"`
	Cache_memory     float64 `json: "cache_memory"`
	Available_memory int     `json: "available_memory"`
	MB_memory        int     `json: "mb_memory"`
}

type Process struct {
	Pid   int      `json: "pid"`
	Name  string   `json: "name"`
	User  int      `json: "user"`
	State int      `json: "state"`
	Ram   int      `json: "ram"`
	Child []Childs `json: "child"`
}

type Cpu struct {
	Processes []Process `json: "processes"`
	Running   int       `json: "running"`
	Sleeping  int       `json: "sleeping"`
	Zombie    int       `json: "zombie"`
	Stopped   int       `json: "stopped"`
	Total     int       `json: "total"`
	Usage     float64   `json: "usage"`
}

type ProcessSend struct {
	Pid   int      `json: "pid"`
	Name  string   `json: "name"`
	User  string   `json: "user"`
	State int      `json: "state"`
	Ram   int      `json: "ram"`
	Child []Childs `json: "child"`
}

type Childs struct {
	Pid  int    `json: "pid"`
	Name string `json: "name"`
}

type CpuSend struct {
	Processes []ProcessSend `json: "processes"`
	Running   int           `json: "running"`
	Sleeping  int           `json: "sleeping"`
	Zombie    int           `json: "zombie"`
	Stopped   int           `json: "stopped"`
	Total     int           `json: "total"`
	Usage     float64       `json: "usage"`
}

// aca empieza la conexion a la base de datos
// Estructura para almacenar los datos de RAM_HISTORIAL
type RamHistorial struct {
	ID  int
	Uso int
}

// Función para insertar datos en la tabla RAM_HISTORIAL
func insertarDatos(db *sql.DB, uso int) error {
	_, err := db.Exec("INSERT INTO RAM_HISTORIAL (uso) VALUES (?)", uso)
	return err
}

// Función para insertar datos en la tabla RAM_HISTORIAL
func insertarDatos2(db *sql.DB, uso int) error {
	_, err := db.Exec("INSERT INTO CPU_HISTORIAL (uso) VALUES (?)", uso)
	return err
}

// Función para obtener una lista de datos de la tabla RAM_HISTORIAL
func obtenerDatos(db *sql.DB) ([]int, error) {
	// Realiza la consulta a la base de datos
	rows, err := db.Query("SELECT uso FROM RAM_HISTORIAL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Crea una lista para almacenar los resultados
	var historial []int

	// Itera sobre los resultados y los agrega a la lista
	for rows.Next() {
		var uso int
		if err := rows.Scan(&uso); err != nil {
			return nil, err
		}
		historial = append(historial, uso)
	}

	// Maneja cualquier error que pueda ocurrir durante la iteración
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return historial, nil
}

// Función para obtener una lista de datos de la tabla CPU_HISTORIAL
func obtenerDatos2(db *sql.DB) ([]int, error) {
	// Realiza la consulta a la base de datos
	rows, err := db.Query("SELECT uso FROM CPU_HISTORIAL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Crea una lista para almacenar los resultados
	var historial []int

	// Itera sobre los resultados y los agrega a la lista
	for rows.Next() {
		var uso int
		if err := rows.Scan(&uso); err != nil {
			return nil, err
		}
		historial = append(historial, uso)
	}

	// Maneja cualquier error que pueda ocurrir durante la iteración
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return historial, nil
}

// Función para manejar las solicitudes HTTP y devolver los datos de uso de RAM en formato JSON
func DataRamHistorial(w http.ResponseWriter, r *http.Request) {
	// Conexión a la base de datos MySQL
	dataSourceName := "root:JAVS27@tcp(localhost:3306)/Proyecto1SO"

	// Abrir una conexión a la base de datos
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Obtener los datos de la tabla RAM_HISTORIAL
	historial, err := obtenerDatos(db)
	if err != nil {
		panic(err.Error())
	}

	// Convertir la lista de usos de RAM en formato JSON
	jsonData, err := json.Marshal(historial)
	if err != nil {
		log.Println("Error al convertir los datos de la RAM a JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Escribir los datos JSON en la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	log.Println("Datos del Historial de la RAM enviados correctamente.")
}

// Función para manejar las solicitudes HTTP y devolver los datos de uso de CPU en formato JSON
func DataCpuHistorial(w http.ResponseWriter, r *http.Request) {
	// Conexión a la base de datos MySQL
	dataSourceName := "root:JAVS27@tcp(localhost:3306)/Proyecto1SO"

	// Abrir una conexión a la base de datos
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Obtener los datos de la tabla RAM_HISTORIAL
	historial, err := obtenerDatos2(db)
	if err != nil {
		panic(err.Error())
	}

	// Convertir la lista de usos de RAM en formato JSON
	jsonData, err := json.Marshal(historial)
	if err != nil {
		log.Println("Error al convertir los datos de la CPU a JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Escribir los datos JSON en la respuesta HTTP
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	log.Println("Datos del Historial de la CPU enviados correctamente.")
}
