import React, { Component } from 'react';
import LineChart from './LineChart';

export class Historial extends Component {
  constructor(props) {
    super(props);
    this.state = {
      memoria: [], // Inicializa el estado con un arreglo vacío
      cpu: [], // Inicializa el estado con un arreglo vacío
    };
  }

  componentDidMount2() {
    // Aquí puedes realizar la solicitud al API y guardar los datos en el estado
    // Ejemplo de solicitud usando fetch:
    fetch('http://localhost:8080/dataramhistorial')
      .then((response) => response.json())
      .then((data) => {
        // Cuando se reciban los datos del API, actualiza el estado con ellos
        this.setState({ memoria: data });
      })
      .catch((error) => {
        console.error('Error al cargar los datos del API:', error);
      });

      fetch('http://localhost:8080/datacpuhistorial')
      .then((response) => response.json())
      .then((data) => {
        // Cuando se reciban los datos del API, actualiza el estado con ellos
        this.setState({ cpu: data });
      })
      .catch((error) => {
        console.error('Error al cargar los datos del API:', error);
      });
  }
  
// Función para cargar datos desde la API
loadDataFromApi(port) {
  fetch(`http://localhost:8080/dataramhistorial`)
    .then((response) => response.json())
    .then((data) => {
      //console.log('RAM data from API:', data);
      this.setState({ memoria: data });
    })
    .catch((error) => {
      console.error('Error fetching RAM data:', error);
    });

  fetch(`http://localhost:8080/datacpuhistorial`)
    .then((response) => response.json())
    .then((data) => {
      //console.log('CPU data from API:', data);
      this.setState({ cpu: data });
    })
    .catch((error) => {
      console.error('Error fetching CPU data:', error);
    });
}

// Maneja el clic en el botón para seleccionar el puerto 3000
handleHistori1ButtonClick() {
  this.setState({ selectedPort: 1 });
  this.loadDataFromApi(1);
}

// Maneja el clic en el botón para seleccionar el puerto 4000
handleHistori2ButtonClick() {
  this.setState({ selectedPort: 2 });
  this.loadDataFromApi(2);
}

componentDidMount() {
  // Inicialmente, carga datos desde el puerto 3000
  this.loadDataFromApi(3000);
}

  render() {
    const { memoria, cpu } = this.state;

    return (
        <>
      



        <h1>Historial del CPU</h1> 
      <div className="center">
        <LineChart dias={[1, 2, 3, 4, 5]} usoRam={cpu} />
      </div>
        <h1>Historial de la RAM</h1> 
      <div className="center">
        <LineChart dias={[1, 2, 3, 4, 5]} usoRam={memoria} />
      </div>
      </>
    );
  }
}