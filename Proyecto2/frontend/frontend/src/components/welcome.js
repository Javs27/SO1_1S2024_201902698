import React from 'react'

export class Welcome extends React.Component {
    render() {
        return (
            <>
            <div class="container">
            <img class="profile-image" src="https://portal.ingenieria.usac.edu.gt/images/logo_facultad/fiusac_negro.png" alt="Imagen de perfil"/>
            <h1 class="profile-name">Pablo Javier Batz Contreras</h1>
            <p class="course-name">Curso: Sistemas Operativos 1</p>
            <p class="course-name">Proyecto 1</p>
            <img src="https://tech.osteel.me/images/2020/03/04/docker-introduction-01.jpg" alt="Docker Logo" width="50%"/>
            </div>
            </>
        )
    }
}