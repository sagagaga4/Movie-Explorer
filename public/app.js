import { HomePage } from "./components/HomePage.js";
import { API } from "./services/API.js"
import './components/AnimatedLoading.js'
import { MovieDetailsPage } from "./components/MovieDetailsPage.js";
import './components/YoutubeEmbed.js'
import { Router } from "./services/Router.js"; 
import { MoviesPage} from "./components/MoviesPage.js";
import { RegisterPage } from "./components/RegisterPage.js";
import { LoginPage } from "./components/LoginPage.js";
import Store from "./services/Store.js";

window.addEventListener("DOMContentLoaded", event =>    {
    app.Router.init();
});

window.addEventListener('error', (event) => {
    console.error('Global error:', event.error);
});

window.app ={
    Router,
    Store,
    showError: (message = "Something went wrong",goToHome=false) => {
        document.getElementById("alert-modal").showModal();
        document.querySelector("#alert-modal p").textContent = message;
        if (goToHome) app.Router.go("/");
    },
    closeError: () => {
        document.getElementById("alert-modal").close()
    },
    search: (event) =>{
        event.preventDefault();
        const q = document.querySelector("input[type=search]").value;
        app.Router.go("/movies?q=" + q);
    
    },
    searchOrderChange: (order) => {
        const urlParams = new URLSearchParams(window.location.search);
        const q = urlParams.get("q");
        const genre = urlParams.get("genre") ?? "";
        app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
    },
    searchFilterChange: (genre) => {
        const urlParams = new URLSearchParams(window.location.search);
        const q = urlParams.get("q");
        const order = urlParams.get("order") ?? "";
        app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
    },
    register: async (event) => {
        event.preventDefault();
        let errors = [];
        const name = document.getElementById("register-name").value;
        const email = document.getElementById("register-email").value;
        const password = document.getElementById("register-password").value
        const passwordConfirm = document.getElementById("register-password-confirm").value;
        //Error handeling 
        if(name.length < 4) errors.push("Enter your complete name.");
        if(email.length < 7) errors.push("Enter a valid email.");
        if(password.length < 4) errors.push("Enter a valid password.");
        if(password != passwordConfirm) errors.push("Passwords do not match (tpesh).");
        if(errors.length == 0){
            const response = await API.register(name,email,password);
            if(response.success){
                app.Store.jwt = response.token;
                app.Router.go("/accounts/")
            } else {
                app.showError(response.message, false);
            }        
        } else {
            app.showError(errors.join(". "), false);
        }
    },
    login: async (event) => {  
        event.preventDefault();
        let errors = []; 
        const email = document.getElementById("login-email").value;
        const password = document.getElementById("login-password").value
        //Error handeling
        if(email.length < 7) errors.push("Enter a valid email.");
        if(password.length < 4) errors.push("Enter a valid password.");
        if(errors.length == 0){
            const response = await API.authenticate(email,password);
            if(response.success){
                app.Store.jwt = response.token;
                app.Router.go("/account/")
            } else {
                app.showError(response.message, false);
            }
        } else {
            app.showError(errors.join(". "), false);
        }
    },
    //Removing user jwt storage and returning to home address
    logout: () => {
        Store.jwt = null;
        app.Router.go("/");
    },
    api : API
}