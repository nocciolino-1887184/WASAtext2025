<script>
export default {
    name: "ConversationInterface",
    props: {
      visible: Boolean,
      messages: Array,
      reactions: Object,
      convos: Array,
      allUsers: Array,
    myId: {
      type: Number,
      required: true
    },
    info:  {
      type: Object,
      default: () => ({
        id: 0,
        avatar: '/default-profile.jpg',
        name: 'Utente sconosciuto',
        tipo: 'privato'        
      }) 
    },
    
    },
    data() {
        return {
            newMess: "",
            idMess: null,
            idInoltro: null,
            convoInolt: "",
            preview: null,
            myNick: sessionStorage.getItem('nickname'),
            selectedFile: null,
            faccine: ['👍', '👎', '😂', '😠'],
            azioni: ['Inoltra', 'Elimina', 'Rispondi'],
            members: [],
            memModal: false,
            isRisp: null,
            idAzione: null
        };
    },
    methods: {
    // funzione unica per mandare foto e testo
    sendCombined() {
      if (this.newMess.trim() === "" && !this.selectedFile) {
        return; // niente da inviare
      }

      this.$emit("send-combined", 
        this.info.id,
        this.newMess || null,
        this.selectedFile || null,
        this.isRisp
      );

      // reset
      this.newMess = "";
      this.isRisp = null;
      this.preview = null;    // visibile solo quando seleziono foto
      this.selectedFile = null;     
    },
  
  allega() {    /* questo serve per non visualizzare altro tasto */
    this.$refs.fileInput.click(); // apre il selettore file
   },
// per mostrare la foto senza inviarla
  handleMess(event) {
    const file = event.target.files[0];
    if (!file) return;

    this.selectedFile = file;
    this.preview = URL.createObjectURL(file);
  },


    selettore(messageId) {  
      if (this.idMess === messageId) {
        this.idMess = null; // chiudi se è già aperto
    } else {
        this.idMess = messageId;
    }
      },
    selettoreA(messageId) {  
       this.idAzione = this.idAzione === messageId ? null : messageId;
     },

    insertEmoji(emoji, messageId) {
      this.$emit('check-react', {emoji, messageId}); 
      this.idMess = null;   
    },
      /* utile per distinguere la comunicazione del backend e la visualizzazione dei dati */
    doAction(azione, messageId) {
      const messaggio = { ...this.messages.find(m => Number(m.id) === Number(messageId))};

      if (!messaggio) {
        console.warn("ERRORE: messaggio non trovato:", messageId);
        return;
      }     
      /* messages serve per la sincronia tra padre  figlio */
          if (azione === 'Inoltra') {
            this.idInoltro = messageId; 
          } else if (azione === 'Elimina') {
            this.deleteMessage(messageId);
          } else if (azione === 'Rispondi') {
            this.isRisp = messaggio;
          }
          this.idMess = null; // Chiude il menu dopo l’azione
    },

    async deleteMessage(messageId) {
      this.$emit('delete-mess', messageId);
      this.idMess = null;      
    },
    
    async forwardMessage() {
       if (!this.idInoltro || !this.convoInolt) {
        alert("Seleziona una destinazione valida.");
        return;
    }

    try {
        // 1. Scomponiamo la selezione (es: "id:15" o "nick:mario")
        const [tipo, valore] = this.convoInolt.split(':'); 

        // 2. Prepariamo il payload ESATTAMENTE come lo aspetta la struct 'Inoltro' del Go
        const payload = {
            messId: this.idInoltro,
            convoId: tipo === 'id' ? parseInt(valore) : 0, // Se è ID, lo mettiamo qui
            dest: tipo === 'nick' ? valore : ""            // Se è Nick, lo mettiamo qui
        };

        const nickname = sessionStorage.getItem('nickname');

        // 3. L'URL: Dato che il Go ignora il parametro nell'URL (è commentato nel tuo codice), 
        // usiamo un valore fittizio o il valore stesso, l'importante è che colpisca la rotta corretta.
        // Se il tuo router Go è impostato su "/conversations/:convoId/forwarded":
        const url = `/conversations/${valore}/forwarded`;

        const response = await this.$axios.post(url, payload, {
            headers: { Authorization: "Bearer " + nickname }
        });

        // 4. Gestione Risposta
        // Il backend torna { "newMess": id, "forwarded": true }
        this.$emit('message-forwarded', {
            target: valore,
            data: response.data
        });

        alert("Messaggio inoltrato con successo!");
        
        // Reset e chiusura modale
        this.convoInolt = ""; 
        this.idInoltro = null;
        this.idAzione = null;
       } catch (error) {
          console.error('Errore durante l’inoltro:', error.message);
        }
    },

    moduloFile() { if (this.info.tipo === 'gruppo' && this.$refs.groupFile)
      {  this.$refs.groupFile.click(); } },

    async changeGp(event){
      const file = event.target.files[0];
      if (!file) return;
      this.$emit('update-avatar', {convId: this.info.id, avatar: file});           
    },

    changeName(event){
      const nuovoNome = prompt("Inserisci il nuovo nome per il gruppo:", this.info.name);
    /* non cambiare props */
      if (nuovoNome && nuovoNome.trim() !== "") {
        // Emetti al genitore
        this.$emit('rename-chat', {convId: this.info.id, nome: nuovoNome});
      }
    },

// qui resituiamo la lista di utenti partecipanti
    async infoMembers() {
      this.memModal = true;
      this.members = [];
      try{
        const response = await this.$axios.get(`/groups/${this.info.id}/members`);
        this.members = response.data;

      }catch(error){
        console.error("Errore check membri:", error);
            alert("Impossibile caricare i membri del gruppo.");
            this.memModal = false;

      }

    },
    closeMem() {
      this.memModal = false;
    },

    addMember(event) {
      const mem = prompt("Che utente vuoi aggiungere nel gruppo?", "");
      this.$emit('new-member', {convId: this.info.id, user: mem});
    },

    async leaveG() {
      try {
        const response = await this.$axios.delete(`/groups/${this.info.id}/members/${this.myId}`);
          alert("Gruppo abbandonato");
        }catch (error) {
          this.errormsg = "Errore nell'uscita dal gruppo : " + error.message;
        }
    },

    formatTime(aString) {
      if (!aString) return '';
    // crea la data dal formato ISO
      const date = new Date(aString);
      if (isNaN(date.getTime())) return 'Data non valida';

      // ora e minuti in formato 24h locale italiana
      return date.toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit', timeZone: 'Europe/Rome' });
    }
  }
};
</script>

<template>

<div class="convo" v-show="visible">
    <div class="intestazione">
      <img :src="info.avatar"
      @error="e => e.target.src = '/default-profile.jpg'"
      class="avatar" alt="avatar" 
      @click="moduloFile" />

      <input
        v-if="info.tipo === 'gruppo'"
        type="file"
        ref="groupFile"
        accept="image/*"
        @change="changeGp"
        style="display: none"
      />

      <div class="chat-name" v-if="this.info.tipo === 'gruppo'" @click="changeName" style="cursor: pointer;">{{ info.name }}</div>
      <div v-if="info.tipo !== 'gruppo'" class="chat-name">
        {{ info.name }}
      </div>
      <!-- non specificare la differenza tra gruppo e chat porta ad errori -->

      <div class="buttons" v-if="this.info.tipo === 'gruppo'">
        <button class="but-group" @click="infoMembers"> <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#info"/> </svg> </button>
        
        <button class="but-group" @click="addMember"> <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus"/></svg> </button>
        <button class="but-group" @click="leaveG"> <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg> </button>        
      </div>

    </div>

<!-- Così metto in un contenitore tutti i dati -->
<div class="chat-messages">  
  <div v-for="mess in messages"
  :key="mess.id"
  :class="['message', Number(mess.sender) === Number(myId) ? 'sender' : 'receiver']">
  <!-- giusto mettere number per non avere problemi con l'id - risolve layout -->

  <div v-if="mess.forwarded" class="forwarded-label">
        <span>Inoltrato</span>
    </div>
  
  <!-- Layout per mostrare le risposte -->
  <div v-if="mess.reply" class="reply-preview">
        <div class="reply-bar"></div> <div class="reply-content">
            <span class="reply-sender">
                {{ mess.replyUser || 'Utente' }} :
            </span>
            <span class="reply-text-short">
                {{ mess.replyText || 'Messaggio originale non disponibile' }}
            </span>
        </div>
    </div>

  <div class="message-meta">  <!-- così visualizzo come si chiama il mandante -->
    <span class="message-sender"> {{  Number(mess.sender) === Number(myId) ? ' ' : mess.sender_name}} </span>    
    <span class="message-time">{{ formatTime(mess.tempo) }} 
      <svg v-if="Number(mess.sender) === Number(myId)" class="feather" :class="{ 'is-read': mess.letto }">
        <use href="/feather-sprite-v4.29.0.svg#eye"/></svg>  </span>
    <button class="emoji" @click="selettore(mess.id)">  
      <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#smile"/></svg>
    </button>
    <button class="actions" @click="selettoreA(mess.id)">
      <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#more-horizontal"/></svg>
    </button>

    <div  v-if="idMess === mess.id"  class="emoji-popup">
        <span v-for="emoji in faccine" :key="emoji" @click="insertEmoji(emoji, mess.id)">
          {{ emoji }}
        </span>
      </div>
   <div  v-if="idAzione === mess.id"  class="action-popup">
        <span v-for="az in azioni" :key="az"
        @click="doAction(az,mess.id)"
        v-show="az !== 'Elimina' || Number(mess.sender) === Number(myId)">
          {{ az }}
        </span>
      </div>

</div>

  <!-- Mostra il testo, se presente -->
  <div class="message-text" :class="{ deleted: mess.deleted }">
    <span v-if="!mess.deleted">{{ mess.text }}</span>
    <span v-else><i>Messaggio eliminato</i></span>
    
  </div>
 
  <!-- Per visualizzare effettivamente la foto -->
   <div v-if="mess.content" class="message-media-wrapper">
    <img
      :src="'data:' + mess.file_type + ';base64,' + mess.content"
      class="message-media"
      alt="immagine"
    />
  </div>

  <!-- Mostra reazioni, se presenti (mess.id è numero) -->
  <div class="reacted"  :class="Number(mess.sender) === Number(myId) ? 'sender' : 'receiver'"
    v-if="reactions[String(mess.id)] && reactions[String(mess.id)].length">
        <span v-for="re in reactions[String(mess.id)]" :key="re.emoji + re.user" :title="re.user">
          {{ re.emoji }} 
          </span>
      </div>


</div>
</div>

 <div v-if="idInoltro" class="forward">
  <div class="forward-box">
    <h3>Seleziona a chi inoltrare il messaggio </h3>
      
    <select v-model="convoInolt">
      <option disabled value="">-select-</option>

      <!-- con questi filtri vedo tutti gli utenti esistenti tranne la convo dove sto -->
        
        <option v-for="c in convos.filter(x => x.id !== info.id && x.id !== myId)"
            :key="'c'+c.id" :value="'id:' + c.id">
            {{ c.name }}
        </option>
    

        <option v-for="u in allUsers.filter(nick => nick !== this.myNick && !convos.some(c => c.name === nick))"
            :key="u" :value="'nick:' + u">
            {{ u }}
        </option>
<!-- chiave con la U per evitare ripetizioni ('u' + u)
    nick serve a scorrere i nomi nelle convo SENZA DOPPIONI -->
    </select>

  <button :disabled="!convoInolt" @click="forwardMessage">
    Inoltra
  </button>
  <button @click="idInoltro = null; convoInolt = ''" class="close-button">Annulla</button>  
  </div>
</div>


<div v-if="memModal" class="modal-overlay" >
  <div class="modal-box">

    <div class="modal-header">
    <h3>  Membri del gruppo
      <span>({{ members.length }})</span>
    </h3>
    <button class="close-icon" @click="closeMem">x</button>
    </div>

    <ul class="members-list">
      <li v-if="!members || members.length === 0" style="text-align:center; color:#888;">
         Nessun membro trovato
      </li>

      <li v-for="(name, index) in members" :key="index">
        <span class="member-icon"> <svg class="feather"> <use href="/feather-sprite-v4.29.0.svg#user"/> </svg> </span> 
        <span>{{ name }}</span> 
      </li>
    </ul>
  </div>
</div>


  <!-- visualizza risposta -->
  <div v-if="isRisp" class="reply-box">
    <div class="reply-text">
      Stai rispondendo a: {{ isRisp.text }}
    </div>
  </div>

   <div v-if="preview" class="preview-container">
      <img :src="preview" class="preview-image"/>
    </div>

    <!-- Barra invio -->
    <div class="message-bar">
      <input
        v-model="newMess"
        type="text"
        placeholder="Scrivi un messaggio..."
      />
      <button @click="allega"> <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#paperclip"/></svg> </button>
      <input
      type="file"
      ref="fileInput"
      accept="image/*"
      @change="handleMess"
      style="display: none"
      />
      <button @click="sendCombined"> Invia </button> 

    </div>
  </div>
</template>

<style scoped>
.convo {
  position: fixed;
  top: 0;
  left: 520px; /* accanto alla sidebar */
  width: calc(100vw - 235px);
  height: 100vh;
  background-color: #d4fcd4; /* verde chiaro */
  display: flex;
  flex-direction: column;
  
}

.intestazione {
  display: flex;
  align-items: center;
  margin-left: 100px;
  margin-right: 285px; /* delimita la destra */
  margin-top: 40px;
  padding: 16px;
  border-bottom: 1px solid #ccc;
  background-color: white;
}

.buttons {
  margin-left: auto;
  display: flex;
  gap: 30px;
}

.but-group {
  background: none;
  border: none;
  font-size: 30px;
  cursor: pointer;
}
/* per bottoni un po' più grandi */
.but-group svg.feather {
  width: 20px;  
  height: 20px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 1rem;
}
.avatar.clickable {
  cursor: pointer;
  transition: transform 0.2s;
}

.chat-name {
  font-size: 1.2rem;
  font-weight: bold;
}

/* Area messaggi */
.chat-messages {
  flex: 1;
  margin-left: 50px;
  padding: 1rem 2rem;
  overflow-y: scroll;
  display: flex;
  flex-direction: column-reverse;
  gap: 2rem;
   /* per far scorrere da sotto */
}

.message {
  position: relative;
  margin-top: 20px;
  padding-bottom: 20px;
  padding-right: 80px;
  padding-left: 10px;
  border-radius: 8px;
  max-width: 40%; 
  bottom: 20px; 
  align-self: flex-start;
  word-wrap: break-word;
  font-size: 1rem;
  line-height: 1.3;
}
.message-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: flex-start; /* default per receiver */
}

/* Contenitore generale della risposta */
.reply-box {
  background: rgba(0, 0, 0, 0.06);
  margin-left: 110px;
  margin-bottom: 6px; 
  border-radius: 6px;
}

.reply-text {
  font-size: 0.80rem;
  font-style: italic;
  color: #444;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.reply-preview {
    background-color: rgba(0, 0, 0, 0.05); /* Sfondo leggermente scuro */
    border-radius: 5px;
    padding: 5px 8px;
    margin-bottom: 5px;
    display: flex;
    font-size: 0.85em;
    cursor: pointer; /* Opzionale: per cliccare e andare al mess */
    border-left: 4px solid #3498db; /* Barra blu a sinistra */
}


.reply-content {
    display: flex;
    flex-direction: column;
}

.reply-sender {
    font-weight: bold;
    color: #555;
    font-size: 0.9em;
}

.reply-text-short {
    color: #777;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis; /* ... se troppo lungo */
    max-width: 200px;
}

.emoji {
  position: absolute;
  bottom: -20px;  /* lo metti sotto il balloon */
  left: 0;        /* bordo sinistro */
  background: white;
  z-index: 10; 
  border: 1px solid #ccc;
  border-radius: 50%;
  padding: 6px;
  cursor: pointer;
  z-index: 2;
  transition: background 0.2s ease;
}

.message-media-wrapper {
  margin-top: 0.5rem;
  display: block;          
  width: auto;   
  
}

.message-media {
  position: relative;
  max-width: 400px;
  max-height: 400px;
  width: auto;
  height: auto;
  display: block;
  border-radius: 8px;
  object-fit: contain;
}

.message.sender {
  align-self: flex-end;
  background-color: #b0eac0;
  text-align: right;
  margin-right: 300px;
}

.message.receiver {
  align-self: flex-start;
  background-color: #f1f1f1;
  text-align: left;
  margin-left: 30px;
}
/* per avere emoji a dx */
.message.receiver .emoji {
  right: 0;
  left: auto;
}

.message-bar {
  display: flex;
  padding: 1rem 100px;
  border-top: 1px solid #ccc;
  background-color: white;
  gap: 10px; /* così i bottoni sono separati */
  max-width: 85%; /* ampiezza rispetto la finestra */
}

.message-bar input {
  flex: 1;
  padding: 0.5rem;
  border: 1px solid #ccc;
  margin-left: 100px;
  border-radius: 6px;
  margin-right: 0.5rem;
  max-width: 80%;
}

.message-bar button {
  padding: 0.5rem 1rem;
  background-color: #2e4f2e;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

.message-bar button:hover {
  background-color: #3e6e3e;
}

.message-meta {
  font-size: 0.9rem;
  color: #555;
  margin-bottom: 4px;
  
}
.message-text.deleted {
  font-style: italic;
  background-color: #b0eac0;
  padding: 6px 10px;
  border-radius: 6px;
}

.message-time {
  position: absolute;
  padding-right: 10px;
  padding-bottom: 10px;
  bottom: 0;
  right: 0;
}

.emoji-popup {
  position: absolute;
  bottom: 10px;
  background: white;
  transform: translateY(0);
  border-radius: 30px; 
  display: flex;
  max-width: 200px;           
  white-space: nowrap; 
  overflow-x: auto;          
  box-sizing: border-box;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  padding: 6px 10px;
  z-index: 999;
  
}
.emoji-popup span {
  font-size: 20px; 
  cursor: pointer;
  padding: 2px 8px;
}
/* per mettere di vedere i box */
.message.receiver .emoji-popup { right: 0;       
  margin-right: -160px;  
}
.message.sender .emoji-popup {  left: 0;
  margin-left: -160px;
}

.preview-container {
  margin-left: 110px;
  width: 150px;
  height: 150px;
  margin-bottom: 10px;
  
}
.preview-image {
  width: 150px;
  height: 150px;
  object-fit: cover;    /* taglia l’immagine senza deformarla */
  border-radius: 4px;
}

.reacted {
  position: absolute;
  bottom: -28px;   /* sotto il messaggio */
  left: 10%;
  background: white;
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 4px 8px;
  display: flex;
  gap: 6px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.15);
  font-size: 15px;
  user-select: none;
  min-width: 40px;
  z-index: 999;
  justify-content: center;
}
.reacted span {
  cursor: default;
  user-select: none;
}
.reacted.sender {  right: 10%; left:auto; }
.actions {
  position: absolute;
  top: 50%; /* così sta a metà messaggio */
  background: transparent;
  border: none;
  padding: 6px;
  cursor: pointer;
  right: -35px;
  transform: translateY(-50%);
  transition: background 0.2s ease;
}
.message.sender .actions { left: -35px; right:auto;}
.action-popup {
  position: absolute;
  top: 50%;
  transform: translateY(-50%); /* per non farlo andare sotto */
  
  background: white;
  border: 1px solid #ccc;
  cursor: pointer;
  border-radius: 8px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.message.receiver .action-popup {
  left: 100%;       
  margin-left: 40px; 
  right: auto;
}

.message.sender .action-popup {
  left: auto;
  right: 100%;
  margin-right: 40px;
}

.forward {
  position: fixed;
  top: 0; left: 0;
  width: 100vw; height: 100vh;
  background-color: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 9999;
}

.forward-box {
  background: white;
  padding: 2rem;
  border-radius: 10px;
  min-width: 300px;
  max-width: 90vw; max-height: 90vh;
  display: flex; flex-direction: column; align-items: center;
  position: relative;
}
.forward-box button {
  margin-top: 1rem;
}

.forwarded-label {
  font-size: 0.90rem;
  color: rgba(0, 0, 0, 0.4);
  font-style: italic;
  margin-top: 4px;
  display: block;
}

.close-button {
  background: none;
  border: none;
  color: #888;
  font-size: 0.9rem;
  margin-top: 1rem;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  top: 0; left: 0;
  width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

/* La finestra bianca */
.modal-box {
  background: white;
  width: 300px;     
  max-height: 70vh;    
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
  display: flex;
  flex-direction: column;
}

/* per avere piccolo bottone x */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-size: 10px;
}

/* Lista pulita */
.members-list {
  list-style: none;
  padding: 0;
  margin: 0;
  overflow-y: auto;
}

.members-list li {
  display: flex;
  align-items: center;
  padding: 8px 0; /* Solo spazio verticale, niente righe */
  gap: 12px;
}

.close-icon {
  background: none; 
  border: none;
  font-size: 20px;
  }

.is-read {
    stroke: #f44336 ; 
}

</style>