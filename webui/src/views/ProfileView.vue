<script>
export default {
  data() {
    return {
      nickname: '',
      tempNick: '', 
      propic: '', 
      errormsg: null,
      showBox: false,
      anteprima: false,
      searching: '',
      output: null
    };
  },
  methods:{
    async modificaNick() {
      const id = sessionStorage.getItem("id");
    /* check perché non voglio avere il nickname vuoto */
      if (!this.tempNick || this.tempNick === "") {
        alert("Inserisci un nickname valido!");
        return; 
      }
      try {
        // ricorda di usare finestrella di inserimento
        const payload = {nickname : this.tempNick}
        /* uso una temporanea per non aggiornare subito la grafica */
        const response = await this.$axios.put(`/users/${id}/nickname`,payload);        
        this.nickname = this.tempNick // aggiornato anche visivamente
        sessionStorage.setItem("nickname", this.nickname);
        alert("Nickname modificato con successo: " + this.tempNick);
        this.showBox = false;
        
      } catch (error) {
        alert("Errore nel cambiamento del nick: " + error.response?.data?.error);
      }
    },

    triggerFileInput() { this.$refs.fileInput.click();  },

    async modificaPropic(event){
      const file = event.target.files[0];

      if (!file) return;
      
      const reader = new FileReader(); // conversione in base 64
      reader.onload= async (e) => {
        try {
          const base64Only = e.target.result.split(',')[1];
          const id = sessionStorage.getItem("id");
          const payload = {propic: base64Only};
          console.log("Invio PUT con payload:", payload);
          const response = await this.$axios.put(`/users/${id}/propic`, payload);
          
          this.propic = `data:image/jpeg;base64,${response.data.propic || base64Only}`;
          alert("Foto profilo aggiornata!");
        } catch (err) {
          alert("Errore nell'aggiornamento della propic: " + err.message);
        }
      };
      reader.readAsDataURL(file);
    },
    async searchNick() {
      if (!this.searching || this.searching.trim() === "") return;
      try {
      const response = await this.$axios.get("/users", {
        params: {nickname: this.searching.trim()}
      });

      this.output = response.data; // contiene nick e foto
      this.searching = ""; // resetta il campo di input
      this.anteprima = true;
    } catch (error) {
        alert("Il nickname cercato non esiste");   
        this.anteprima = false;
       }
     },

   },

  mounted() {
    this.nickname = sessionStorage.getItem("nickname") 
    // per quanto riguarda l'immagine profilo, decidiamo di prenderla sempre dal database
    const id = sessionStorage.getItem("id");

    this.$axios.get(`/users/${id}`)
      .then(response => {
        const data = response.data;

        if (data.propic) {
          this.propic = `data:image/jpeg;base64,${data.propic}`;
        } else {
          this.propic = '/default-profile.jpg';  // se non ho foto settata
        }

      })
      .catch(error => {
        console.error("Errore nel recupero delle info utente:", error);
        this.propic = '/default-profile.jpg';
      });
    }

};
</script> 

<template>
            <!-- Scritta simile alla homepage -->
  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h3 text-muted">Profilo</h1>
  </div>

   <h5>Ricerca un utente dal nickname</h5>
   <div class="search-bar"> 
  <input
        v-model="searching"
        type="text"
        placeholder="Cerca un utente"
    /> <button @click="searchNick()">Cerca</button>
  </div>


  <div class="container d-flex flex-column align-items-center justify-content-center" style="height: 80vh;">
    <button @click="triggerFileInput" class="custom mb-3"> <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit"/></svg> </button>
    <img
      :src="propic"
      alt="Image not found"
      class="rounded-circle mb-3"
      style="width: 200px; height: 200px; object-fit: cover;"
    /> 
    <input
      type="file"
      ref="fileInput"
      accept="image/*"
      @change="modificaPropic"
      style="display: none"
    />
    

    <h2>{{ nickname }}
     <button @click="showBox = true" class="custom">
      <svg class="feather">
      <use href="/feather-sprite-v4.29.0.svg#edit" />
        </svg>
    </button>
    </h2> 
    
	<div v-if="showBox" class="box" @click.self="showBox = false">
 	<div class="box-content">
	<button @click="showBox = false" class="close-box-btn">×</button>
    <h3>Modifica il nickname</h3>
    
    <input v-model="tempNick" type="text" placeholder="Nuovo nome" />

    <div class="box-actions">
      <button @click="modificaNick" class="btn btn-success">Crea</button>
      </div>
  </div>
</div>
    <div class="d-flex align-items-center mt-2">
      <span class="status-dot me-2"></span>
      <span class="text-success">Online</span>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>

 

  
<div v-if="anteprima && output" class="preview" @click.self="anteprima = false">
  <div class="preview-box">
    <button class="close-box-btn" @click="anteprima = false">×</button>
    <img :src="output.propic ? 'data:image/png;base64,' + output.propic : '/default-profile.jpg'"
         alt="Foto profilo"
         class="preview-img" />
    <h3>{{ output.nickname }}</h3>
  </div>
</div>

  
</template>


<style scoped>
.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: #28a745; /* verde */
  display: inline-block;
}
.box {
  position: fixed;
  top: 0; left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.box-content {
  position: relative;
  background: white;
  padding: 3rem;
  border-radius: 8px;
  min-width: 300px;
}

.box-content input {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 1rem;
}

.box-actions {
  display: flex;
  justify-content: space-between;
}
.close-box-btn {
  position: absolute;
  top: 10px;
  right: 15px;
  border: none;
  background: none;
  font-size: 24px;
  cursor: pointer;
  color: #000; }

.search-bar {
  top: 155px;
  display: flex;
  background-color: white;
  gap: 10px;
  font-size: 0.9rem;
  border-radius: 6px;
  max-width: 300px;
}

.search-bar input {
  flex: 1;
  padding: 0.3rem 0.5rem;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.search-bar button {
  padding: 0.3rem 0.8rem;
  background-color: #2e4f2e;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.preview {
  position: fixed;
  top: 0; left: 0;
  width: 100vw; height: 100vh;
  background-color: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 9999;
}

.preview-box {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  min-width: 300px;
  max-width: 90vw; max-height: 90vh;
  display: flex; flex-direction: column; align-items: center;
  position: relative;
}

.preview-img {
  width: 120px; height: 120px;
  object-fit: cover; border-radius: 50%;
  margin-bottom: 1rem;
}

.custom {
  border: 2px solid #90ee90; /* verde chiaro */
  background-color: transparent;
  padding: 6px 6px;
  border-radius: 6px;
}

</style>
