import React from 'react';
import PieChart from './PieChart';

export class Memory extends React.Component {
    state = {
        data: [['x', 'Memoria RAM'], [1, 0], [2, 0], [3, 0], [4, 0], [5, 0], [6, 0], [7, 0], [8, 0], [9, 0], [10, 0], [11, 0], [12, 0], [13, 0], [14, 0], [15, 0]],
        memoria: { Total_memory: 0, Free_memory: 0, Used_memory: 0, Available_memory: 0, MB_memory: 0 }
    }

    componentDidMount() {
        this.loadDataFromApi('http://localhost:8080/dataram');
    }

    loadDataFromApi(apiUrl) {
        fetch(apiUrl)
            .then((response) => response.json())
            .then((dataFromServer) => {
                console.log('Data from API:', dataFromServer);
                this.fillData();
                this.setState({ memoria: dataFromServer });
            })
            .catch((error) => {
                console.error('Error fetching data:', error);
            });
    }

    render() {
        return (
            <>
                <div className="center">
                    <PieChart porcentaje={this.state.memoria.Used_memory} />
                </div>

                <h2>Uso de la RAM: {this.state.memoria.Used_memory}%</h2>

                <div className='row'>
                    <div className='col'>
                        <br />
                        <br />
                        <p>
                            Memoria total: {this.state.memoria.Total_memory} MB
                        </p>
                        <p>
                            Uso de memoria: {this.state.memoria.Used_memory}%
                        </p>
                        <p>
                            Memoria disponible: {100 - this.state.memoria.Used_memory}%
                        </p>
                    </div>
                </div>
            </>
        );
    }

    fillData() {
        const encabezado = ['x', 'Memoria RAM'];
        const inputData = [Number(this.state.data[15][0]) + 1, this.state.memoria.Used_memory];
        const datos = [encabezado, ...this.state.data.slice(1, -1), inputData];
        this.setState({ data: datos });
    }
}
