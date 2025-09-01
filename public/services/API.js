export const API = {
    baseURL: '/api/',
    getTopMovies: async () => {
        return await API.fetch("movies/top");
    },
    getRandomMovies: async () => {
        return await API.fetch("movies/random");
    },
    getMovieById: async (id) => {
        return await API.fetch(`/movies/${id}`);
    },
    searchMovies: async (q, order, genre) => {
        return await API.fetch(`/movies/search`, {q, order, genre})
    },
    getGenres: async () => {
        return await API.fetch("genres/");
    },
    loadGenres: async () => {
        const genres = await API.getGenres("genres");
    },
    register: async (name, email, password) => {
        return await API.send("account/register/",{name,email,password});
    },
    authenticate: async (email, password) => {
        return await API.send("account/authenticate/", {email, password})
    },   
    send: async (ServiceName, data) => {
        try{
            const response = await fetch(API.baseURL + ServiceName, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });
            const result = await response.json();
            return result;
            }catch(e){
            console.error(e);
            }
    },
    fetch: async (service, args) => {
        try {
            const queryString = args ? new URLSearchParams(args).toString() : "";
            const response = await fetch(API.baseURL + service + '?' + queryString);
            const result = await response.json();
            return result;
        } catch (e) {
            console.error(e);
            app.showError();
        }
    }
}

export default API;
