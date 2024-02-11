import React, { useState, useRef, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [photos, setPhotos] = useState([]);
  const videoRef = useRef(null);
  const mediaStreamRef = useRef(null);
  const isPlayingRef = useRef(false);

  useEffect(() => {
    const startCamera = async () => {
      try {
        const stream = await navigator.mediaDevices.getUserMedia({ video: true });
        mediaStreamRef.current = stream;
        if (videoRef.current && !isPlayingRef.current) {
          videoRef.current.srcObject = stream;
          videoRef.current.addEventListener('loadedmetadata', () => {
            videoRef.current.play().then(() => {
              isPlayingRef.current = true;
            }).catch(error => {
              console.error('Error playing video:', error);
              alert('Failed to play video. Please try again.');
            });
          });
        }
      } catch (error) {
        console.error('Error starting camera:', error);
        alert('Failed to start camera. Please make sure your camera is connected and try again.');
      }
    };

    startCamera();

    return () => {
      stopCamera();
    };
  }, []);

  const capturePhoto = async () => {
    try {
      const canvas = document.createElement('canvas');
      canvas.width = videoRef.current.videoWidth;
      canvas.height = videoRef.current.videoHeight;
      const context = canvas.getContext('2d');
      context.drawImage(videoRef.current, 0, 0, canvas.width, canvas.height);

      const base64Image = canvas.toDataURL('image/jpeg');
      setPhotos(prevPhotos => [...prevPhotos, base64Image]);
    } catch (error) {
      console.error('Error capturing photo:', error);
      alert('Failed to capture photo. Please try again.');
    }
  };

  const sendPhotosToBackend = async () => {
    try {
      const formData = new FormData();
      photos.forEach((photo, index) => {
        // Convertir la imagen base64 a Blob
        const blob = fetch(photo)
          .then(res => res.blob())
          .then(blob => {
            // Agregar cada foto como un archivo a FormData
            formData.append('photos', blob, `photo-${index}.jpg`);
          });
      });
  
      // Esperar a que todas las fotos sean procesadas y agregadas a formData
      await Promise.all(photos.map((photo, index) => fetch(photo)
        .then(res => res.blob())
        .then(blob => {
          formData.append('photos', blob, `photo-${index}.jpg`);
        })
      ));
  
      // Realizar la petición POST con FormData
      const response = await axios.post('http://localhost:3001/uploadPhotos', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
  
      console.log(response.data);
      setPhotos([]);
    } catch (error) {
      console.error('Error sending photos to backend:', error);
      alert('Failed to send photos to backend. Please try again.');
    }
  };

  const stopCamera = () => {
    if (mediaStreamRef.current) {
      mediaStreamRef.current.getTracks().forEach(track => track.stop());
    }
  };

  const getPhotosFromBackend = async () => {
    try {
      const response = await axios.get('http://localhost:3001/getPhotos');
      const photosFromBackend = response.data;
      // Mapear la ruta relativa a la URL completa del servidor
      photosFromBackend.forEach(photo => {
        // Actualiza la URL de la imagen con la ruta completa del servidor
        photo.image = `http://localhost:3001/${photo.image}`;
      });
  
      // Abre una nueva ventana del navegador para mostrar las fotos
      const newWindow = window.open('');
      newWindow.document.write('<html><head><title>Fotos desde el backend</title></head><body>');
      photosFromBackend.forEach(photo => {
        newWindow.document.write(`<img src="${photo.image}" />`);
      });
      newWindow.document.write('</body></html>');
    } catch (error) {
      console.error('Error getting photos from backend:', error);
      alert('Failed to get photos from backend. Please try again.');
    }
  };
  

  return (
    <div className="App">
      <header className="App-header">
        <div className="camera-container">
          <video ref={videoRef} className="video-element" autoPlay muted></video>
          <button className="capture-button" onClick={capturePhoto}>Take Photo</button>
        </div>
        <div className="photo-gallery">
          {photos.map((photo, index) => (
            <img key={index} src={photo} alt={`Captured ${index}`} className="captured-image" />
          ))}
        </div>
        {photos.length > 0 && (
          <button className="send-button" onClick={sendPhotosToBackend}>Send Photos to Backend</button>
        )}
        <button className="get-photos-button" onClick={getPhotosFromBackend}>Get Photos from Backend</button>
      </header>
    </div>
  );
}

export default App;
