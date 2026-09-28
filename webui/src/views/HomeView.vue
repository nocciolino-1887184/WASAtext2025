<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			nickname: "",
			finestraSelezione: false,
			tipo: "",
			formConv: { nome: "", membri: "" },
			
		};
	},
	methods: {
		/* async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},  */
		
		openModal(type) {
      	 this.tipo = type;
      	 this.finestraSelezione = true;
    	},
		async createConvo() {
			const payload = {		// aggingi altri
				members: this.formConv.membri
					.split(",")
					.map(s => s.trim())
					.filter(n => n !== ""),  // controlla sia numero valido
				name: this.formConv.nome,
				tipo: this.tipo				
		};
		
			try {
				const nik = sessionStorage.getItem('nickname');
				const response = await this.$axios.post("/conversations",payload, {
				headers: { Authorization: "Bearer " + nik }
					}); 
				
				const convoId = response.data.id;
				const id = sessionStorage.getItem("id");
				this.$router.push(`/users/${id}/conversations`);
			} catch (error) {
				alert("Errore nella creazione della conversazione: " + error.response?.data?.error);
			}
    },
	},
	mounted() {
		this.nickname = sessionStorage.getItem('nickname')
		// this.refresh()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h3 text-muted">Homepage</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
			<div class="btn-group me-2">
				<!--<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button> -->
				</div>
		<!--		<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button> 
				</div> -->
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	
	<div class="container-fluid position-relative">

    <div class="position-absolute top-0 end-0 p-3">
    </div>

    <div class="d-flex flex-column align-items-center justify-content-center" style="height: 80vh;">
      <h1 class="display-4">Benvenuto/a, {{ nickname }}!</h1>
      <p class="lead">Cosa vuoi fare?</p>

      <div class="mt-3 d-flex gap-3">
        <button @click="openModal('chat')" class="btn custom-btn">Nuova Chat</button>
        <button @click="openModal('gruppo')" class="btn custom-btn">Nuovo Gruppo</button>
	<div v-if="finestraSelezione" class="modal-backdrop" @click.self="finestraSelezione = false">
 	<div class="modal-content">
	<button @click="finestraSelezione = false" class="close-modal-btn">×</button>
    <h3>Crea {{ tipo === 'gruppo' ? 'un gruppo' : 'una chat' }}</h3>
    
    <input v-if="tipo === 'chat'" v-model="formConv.nome" type="text" placeholder="Destinatario" />
	<input v-if="tipo === 'gruppo'"  v-model="formConv.nome" type="text" placeholder="Nome del gruppo" />
    <input v-if="tipo === 'gruppo'"  v-model="formConv.membri" type="text" placeholder="Partecipanti (separati da virgola)" />

    <div class="modal-actions">
      <button @click="createConvo" class="btn btn-success">Crea</button>
     
    </div>
  </div>
</div>
      </div>
    </div>
  </div>
  </div>
</template>

<style scoped>

.custom-btn {
  border: 2px solid #90ee90; /* verde chiaro */
  background-color: transparent;
  padding: 10px 20px;
  font-weight: 500;
  border-radius: 6px;
  transition: background-color 0.2s ease;  /* così da avere verde se ci passo sopra mouse */
}
.custom-btn:hover {
  background-color: #d4fcd4; 
}

.modal-backdrop {
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

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  width: 400px;
  max-width: 90%;
}

.modal-content input {
  width: 100%;
  padding: 0.5rem;
  margin-bottom: 1rem;
}

.modal-actions {
  display: flex;
  justify-content: space-between;
}
.close-modal-btn {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  position: absolute;
  top: 0.5rem; right: 0.5rem;
  cursor: pointer; }

</style>
