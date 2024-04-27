import React from 'react';
import PieChart from './PieChart';

export class Cpu extends React.Component {
  state = {
    data: [['x', 'Memoria RAM'], [1, 0], [2, 0], [3, 0], [4, 0], [5, 0], [6, 0], [7, 0], [8, 0], [9, 0], [10, 0], [11, 0], [12, 0], [13, 0], [14, 0], [15, 0]],
    cpuData: { Processes: [], Running: 0, Sleeping: 0, Zombie: 0, Stopped: 0, Total: 0, Usage: 0 },
    searchTerm: '', // Para el término de búsqueda
  };

  componentDidMount() {
    this.loadDataFromApi('http://localhost:8080/datacpu');
  }

  loadDataFromApi(apiUrl) {
    fetch(apiUrl)
      .then((response) => response.json())
      .then((dataFromServer) => {
        console.log('Data from API:', dataFromServer);
        this.fillData();
        this.setState({ cpuData: dataFromServer });
      })
      .catch((error) => {
        console.error('Error fetching data:', error);
      });
  }

  handleSearchTermChange = (event) => {
    const searchTerm = event.target.value.toLowerCase();
    this.setState({ searchTerm });

    // Buscar la primera fila que coincide con el término de búsqueda
    const matchingRow = this.state.cpuData.Processes.find(
      (element) => element.Name.toLowerCase() === searchTerm
    );

    // Si se encontró una fila que coincide, desplázate hacia ella
    if (matchingRow) {
      const index = this.state.cpuData.Processes.indexOf(matchingRow);
      const row = document.getElementById(`row-${index-5}`);
      if (row) {
        row.scrollIntoView({ behavior: 'smooth' });
      }
    }
  };

  render() {
    const filteredProcesses = this.state.cpuData.Processes.filter((element) => {
      const processName = element.Name.toLowerCase();
      const searchTerm = this.state.searchTerm.toLowerCase();
      return processName.includes(searchTerm);
    });

    return (
      <>
        <div className='center'>
          <PieChart porcentaje={this.getUsage()}/>  
        </div>
            
        <h2>Uso del CPU: {this.getUsage()}%</h2>

       
        
      </>
    );
  }

  getUsage() {
    return this.state.cpuData.Usage;
  }

  fillData() {
    var encabezado = ['x', 'Uso CPU'];
    var inputData = [Number(this.state.data[15][0]) + 1, this.state.cpuData.Usage];
    console.log(this.state.data[7]);
    var datos = [];
    datos.push(encabezado);
    for (let i = 0; i < 15; i++) {
      if (this.state.data[i + 2]) {
        datos.push(this.state.data[i + 2]);
      }
    }
    datos.push(inputData);
    this.setState({ data: datos });
  }

  killProcess(pid) {
    // Crear un objeto JavaScript con la estructura { "Pid": pid }
    const data = { Pid: pid };
      
    fetch("http://localhost:8080/kill", {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data) // Convertir el objeto a JSON
    }).then(async response => {
      const json = await response.json();
      if (json.value !== false) {
        alert("Se mató el proceso");
        console.log(json);
      } else {
        alert("Error al matar el proceso");
      }
    });
  }
}
