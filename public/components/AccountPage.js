export class AccountPage extends HTMLElement{
    connectedCallback(){
    const template = document.getElementById("template-account");
    const content = template.content.cloneNode(true);
    this.appendChild(content);
    }
}
customElements.define("account-page", AccountPage);
