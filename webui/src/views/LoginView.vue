<script>
export default {
  name: "LoginPage",
  data() {
    return {
      nickname: "",
      loggedIn: false
    };
  },
  methods: {
  async handleLogin() {
    if (this.nickname.trim() === "") {
      alert("Inserisci un nickname valido.");
      return;
    }

    try {
      const response = await this.$axios.post(`/session`, {nickname: this.nickname   });

      sessionStorage.setItem("nickname", this.nickname);
      sessionStorage.setItem("loggedIn", "true")
        // voglio salvarmi l'ID
      const id = response.data.id;
      if (id) {
        sessionStorage.setItem("id", id.toString());
      } else {
        console.warn("Attenzione: userId non ricevuto dal backend.");
      }
      this.$router.push("/dashboard");

    } catch (error) {
      console.error("Errore dettagliato:", error);
      const message = error.response?.data?.message || "Errore di connessione al server";
      alert("Errore: " + message);
    }
  },
},

};
</script>

<template>
  <div class="login-page">
    <h1 class="title">WASA Text</h1>

    <div class="login-box">
      <h2>Login</h2>
      <input
        v-model="nickname"
        type="text"
        placeholder="Inserisci nickname"
        class="input"
      />
      <button @click="handleLogin" class="btn">Register / Login</button>
    </div>
  </div>
</template>



<style scoped>
.login-page {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: #d4fcd4; /* verde chiaro */
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999; /* assicura che stia sopra sidebar etc per non mostrre la barra */
  overflow: hidden;
}

.title {
  font-size: 3rem;
  margin-bottom: 2rem;
  color: #2e4f2e;
}

.login-box {
  background-color: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  min-width: 300px;
}

.input {
  width: 100%;
  padding: 0.75rem;
  margin-top: 1rem;
  margin-bottom: 1rem;
  border: 1px solid #ccc;
  border-radius: 6px;
  font-size: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  background-color: #2e4f2e;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
}

.btn:hover {
  background-color: #3e6e3e;
}

</style>
