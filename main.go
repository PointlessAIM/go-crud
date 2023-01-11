package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func dbConn() (conexion *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "sistema"
	conexion, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var templates = template.Must(template.ParseGlob("template/*"))

// aquí se determinan las URLs junto con su función correspondiente
// también verifica si se está haciendo conexión con el puerto.
func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)


	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)

}

// aquí se determina el modelo de la base de datos
type Empleado struct{
	Id int
	Nombre string
	Correo string
}

// primer view, muestra la información de la base de datos basandose
// en el modelo "Empleado"
func Index(w http.ResponseWriter, r *http.Request) {

	conectionEstablished:=dbConn()

	registros, err := conectionEstablished.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	
	empleado := Empleado{}
	arregloEmpleado:=[]Empleado{}

	//copia los registros de la base de datos y los almacena en
	// arregloEmpleado
	for registros.Next(){
		var id int
		var nombre, correo string
		err= registros.Scan(&id,&nombre,&correo)
		if err!=nil {
			panic(err.Error())
		}
		empleado.Id= id
		empleado.Nombre= nombre
		empleado.Correo= correo

		arregloEmpleado=append(arregloEmpleado, empleado)
	}

	templates.ExecuteTemplate(w, "index", arregloEmpleado)

}

//redirecciona a la plantlla para agregar nueva información
func Crear(w http.ResponseWriter, r *http.Request) {
	
	templates.ExecuteTemplate(w, "crear", nil)

}

//inserta información a la DB
func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST"{
		nombre:= r.FormValue("Nombre")
		correo:= r.FormValue("Correo")

		conectionEstablished:=dbConn()
		insertarRegistros, err := conectionEstablished.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?, ?)")

		if err != nil {
			panic(err.Error())
		}
	
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w,r,"/",http.StatusMovedPermanently)

	}
	//templates.ExecuteTemplate(w, "crear", nil)

}

func Borrar (w http.ResponseWriter, r *http.Request){

	idEmpleado:= r.URL.Query().Get("id")
	fmt.Println(idEmpleado)

	conectionEstablished:=dbConn()
	borrarRegistros, err := conectionEstablished.Prepare("DELETE FROM empleados WHERE id=?")

	if err!=nil{panic(err.Error())}

	borrarRegistros.Exec(idEmpleado)
	http.Redirect(w,r,"/",http.StatusMovedPermanently)

}

func Editar(w http.ResponseWriter, r *http.Request){
	idEmpleado:= r.URL.Query().Get("id")
	log.Println(idEmpleado)

	conectionEstablished:=dbConn()
	registro, err := conectionEstablished.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}

	for registro.Next(){
		var id int
		var nombre, correo string
		err= registro.Scan(&id,&nombre,&correo)
		if err!=nil{
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo 
	}
	log.Println(empleado)
	templates.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request){

	if r.Method == "POST"{
		id:= r.FormValue("id")
		nombre:= r.FormValue("Nombre")
		correo:= r.FormValue("Correo")

		conectionEstablished:=dbConn()
		actualizarRegistros, err := conectionEstablished.Prepare(" UPDATE empleados SET nombre=?, correo=? WHERE id=? ")

		if err != nil {panic(err.Error())}
	
		actualizarRegistros.Exec(nombre, correo, id)
		http.Redirect(w,r,"/",http.StatusMovedPermanently)

	}
}