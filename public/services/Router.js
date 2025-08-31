import { routes } from "./Routes.js";

export const Router = {
    init: () => {
        window.addEventListener("popstate", () => {
            Router.go(location.pathname, false);
        });      
        // Enchance current links in the document
        document.querySelectorAll("a.navlink").forEach(a => {
            a.addEventListener("click", event => {
                event.preventDefault();
                const href = a.getAttribute("href");
                Router.go(href);
            });
        });  

        // Process initial URL   
        Router.go(location.pathname + location.search);
    },
    go: (route, addToHistory=true) => {
        if (addToHistory) {
            history.pushState(null, "", route);
        }
        const routePath = route.includes('?') ? route.split('?')[0] : route;
        let pageElement = null;
        
        let needLogin = false;

        for (const r of routes) {
            if (typeof r.path === "string" && r.path === routePath) {
                pageElement = new r.component();
                break;
            } else if (r.path instanceof RegExp) {
                const match = r.path.exec(route);
                if (match) {
                    const params = match.slice(1);
                    pageElement = new r.component();
                    pageElement.params = params;                    
                    break;
                }
            }
            needLogin = r.loggedIn == true
        }

        if(pageElement){
            if(needLogin && app.Store.loggedIn == false){
                app.Router.go("/account/go")
                return;
            }

        }



        if (pageElement==null) {
            pageElement = document.createElement("h1");
            pageElement.textContent = "Page not found";
        }       
        //Inserting the new page in the UI
        const oldPage = document.querySelector("main").firstElementChild;
        if(oldPage) oldPage.style.viewTransitionName = "old";
        pageElement.style.viewTransitionName = "new";
        
        function updatePage() {
            document.querySelector("main").innerHTML = "";
            document.querySelector("main").appendChild(pageElement);
        }


        if(document.startViewTransition){
            //We don't do a transition 
            updatePage();
        }   else {
            //We do a transition
            document.startViewTransition(() =>{
                updatePage();
            });
        }
    }
}

 