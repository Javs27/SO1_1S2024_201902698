import React, { useEffect, useRef } from 'react';
import Chart from 'chart.js';

const LineChart = ({ dias, usoRam }) => {
  const chartRef = useRef(null);

  useEffect(() => {
    const data = {
      labels: dias, // Etiquetas en el eje X (por ejemplo, ['Día 1', 'Día 2', 'Día 3'])
      datasets: [
        {
          label: 'Uso',
          data: usoRam, // Datos de uso de RAM (por ejemplo, [50, 60, 70])
          borderColor: 'blue', // Color de la línea para el uso de RAM
          borderWidth: 2,
          fill: false,
        },
        // Puedes agregar más conjuntos de datos si es necesario
      ],
    };

    const options = {
      responsive: true,
      scales: {
        x: {
          beginAtZero: true,
          title: {
            display: true,
            text: 'Día',
          },
        },
        y: {
          beginAtZero: true,
          title: {
            display: true,
            text: 'Uso de RAM',
          },
        },
      },
    };

    const ctx = chartRef.current.getContext('2d');
    new Chart(ctx, {
      type: 'line',
      data: data,
      options: options,
    });
  }, [dias, usoRam]);

  return (
    <div className="line-chart">

      <canvas ref={chartRef} width={100} height={100} />
    </div>
  );
};

export default LineChart;
