import React, { useEffect, useRef } from 'react';
import Chart from 'chart.js';

const PieChart = ({ porcentaje }) => {
  const chartRef = useRef(null);
  const libre = 100 - porcentaje; // Calcular el valor "Libre" aquí

  useEffect(() => {
    // Datos y opciones de la gráfica
    const data = {
      labels: ['Uso', 'Libre'],
      datasets: [
        {
          data: [porcentaje, libre], // Utilizar 'libre' en lugar de calcularlo
          backgroundColor: ['red', 'blue'],
        },
      ],
    };

    const options = {
      responsive: true,
    };

    // Crear la gráfica en el elemento canvas
    const ctx = chartRef.current.getContext('2d');
    new Chart(ctx, {
      type: 'pie',
      data: data,
      options: options,
    });
  }, [porcentaje, libre]); // Agregar 'porcentaje' y 'libre' como dependencias del efecto

  return (
    <div className="pie-chart">
      
      <canvas ref={chartRef} width={300} height={300} />
    </div>
  );
};

export default PieChart;

