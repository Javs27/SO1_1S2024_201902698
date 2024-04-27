import React from 'react'
import { Link } from 'react-router-dom'

export class Navbar extends React.Component {
    render() {
        return (
            <div>
                
                <nav class="navbar navbar-expand-lg navbar-light bg-primary">
                    <div class='container'>
                            <a class="navbar-brand" href='/' >
                             MODULOS DEL KERNEL
                            </a>
                        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                            <span class="navbar-toggler-icon"></span>
                        </button>
                        <div class="collapse navbar-collapse" id="navbarSupportedContent">
                            <ul class="navbar-nav mr-auto">
                    
                                <li class="nav-item active">
                                    <Link to='/cpu'>
                                        {/* eslint-disable jsx-a11y/anchor-is-valid */}
                                        <a class="nav-link" >CPU</a>
                                    </Link>
                                </li>
                                <li class="nav-item active">
                                    <Link to='/memory'>
                                        
                                        <a class="nav-link" >RAM</a>
                                    </Link>
                                </li>
                                <li class="nav-item active">
                                    <Link to='/historial'>
                                        <a class="nav-link" >HISTORIAL</a>
                                    </Link>
                                </li>
                            </ul>
                        </div>
                    </div>
                </nav>
            </div>
        )
    }
}
