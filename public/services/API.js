export const API = {
    baseURL: '/api/',
    getJWT: () => {
        return localStorage.getItem('jwt');
    },
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
        return await API.send("account/register/",{name, email, password});
    },
    login: async(email, password) => {
        return await API.send("account/authenticate/", {email, password})
    },
    authenticate: async(email, password) => {
        return await API.send("account/authenticate/", {email, password})
    },
    /*
    authenticate: async (email, password) => {
        return await API.send("account/authenticate/", {email, password})
    },
    */

    getFavorites: async () => {
            return await API.fetch("account/favorites/");
    },     
    
    getWatchlist: async () => {
            return await API.fetch("account/watchlist/");
    },     

    saveToCollection: async (movie_id, collection) => {
        return await API.send("account/save-to-collection/", {
            movie_id, collection
        });
    },

    send: async (serviceName, args) => {
        const jwt = API.getJWT();
        const response = await fetch(API.baseURL + serviceName, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": jwt ? `Bearer ${jwt}` : null
            },
            body: JSON.stringify(args)
        });
        const result = await response.json();
        return result;
    },

    fetch: async (service, args) => {
        try{
        const jwt = API.getJWT();
        const queryString = args ? new URLSearchParams(args).toString() : "";
        const response = await fetch(API.baseURL + service + '?' + queryString, {
            headers: {
                "Authorization": jwt ? `Bearer ${jwt}` : null
            }
        });
        const result = await response.json();
        return result;
    } catch (e) {
        console.error(e);
        throw e;
        }
    }
}

export default API;
