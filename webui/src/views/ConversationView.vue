<script>
import SideConvo from '../components/SideConvo.vue'
import ConvoInterface from '../components/ConvoInterface.vue'

export default {
  components: { SideConvo, ConvoInterface },
	data() {
		return {
    errormsg: null, 
		showConv: true,
    showChat: true,
		chats: [],
    messages: [],
		userId: 0,	
    intervalId: null,
    chatAttiva: null,
    risposta: null,
    numMess: 0,
    allUsers: [],
    reactions: {}
    };
  },
  computed: {   // parte obbligatoria per tenere dati in chache e aggiornare la chat (readOnly)
    interaChat() {
      if (!this.chatAttiva) return null;
      return this.chats.find(c => c.id === this.chatAttiva) || null;
    }
  },
   methods: {
    async showConvos() {   // così prendo dati dal database e li passo alla grafica
      try {
        const id = sessionStorage.getItem('id');
        this.userId = Number(id);
        
        const response = await this.$axios.get(`/users/${id}/conversations`);
        
        this.chats = response.data.map(conv => {
          const isGroup = conv.categoria === "gruppo";
          let finalAvatar = "/default-profile.jpg";
          const raw = conv.icona ? conv.icona.trim() : "";

          if (raw !== "") {
            if (raw.startsWith("data:") || raw.startsWith("http")) {
              finalAvatar = raw;
            } else {
              // Identifichiamo il tipo (Uniformiamo PNG e JPG)
              const type = raw.startsWith("iVBOR") ? "png" : "jpeg";
              finalAvatar = `data:image/${type};base64,${raw}`;
            } }
          const nomeC = isGroup? conv.nome : conv.ricevitore;  // uno dei due è null
        
          var letto = conv.lettura;
          if (conv.id === this.chatAttiva) {
            if (letto){
              letto = null; 
              this.notificaLet(conv.id);
              }
          }
          return {
            id: conv.id,
            name:  nomeC || "Sconosciuto", 
            ricevitore: conv.ricevitore || null,
            avatar: finalAvatar,
            tipo: conv.categoria,
            lastMessage: conv.tipo_mess === "text" ? conv.mess : ( (conv.tipo_mess === "image") ?  "Allegato" : "Ancora nessun messaggio" ) ,
            lettura: letto,
            time: conv.tempo || null
        };
      });
             
      } catch (error) {
        this.errormsg = "Errore nel recupero delle conversazioni: " + error.message;
      }
    },
        // così quando sto già dentro la convo, me lo segnala come letto!!
    async notificaLet(chatId) {
    try {
        const nickname = sessionStorage.getItem('nickname');
          // il body del PUT È VUOTO MA VA SCRITTO
        await this.$axios.put(`/conversations/${chatId}/letti`, {}, {
             headers: { Authorization: "Bearer " + nickname }
        });
    } catch (error) {
        console.error("Errore auto-lettura:", error);
    }
},
    /* richiamo il nickname del mandante */
    async getNickById(userId) {    
      try {
        const res = await this.$axios.get(`/users/${userId}`);
        const nickname = res.data.username;
        return nickname;
      } catch (err) {
        return 'Utente';
      }
    },

    // per fornire id della chat
    async selezionaChat(chat){
      this.chatAttiva = chat.id;
      this.showMess(chat.id);
      const index = this.chats.findIndex(c => c.id === chat.id);
      
      if (index !== -1) {
        const daAggiornare = this.chats[index].lettura; 
        // qui la chat intera si aggiorna da sola
        this.numMess = 0;
                                      // aggiornamento del padre della lettura di nuovi messaggi
      if (daAggiornare) {
        this.chats[index].lettura = null;
        chat.lettura = null;
            
      try {
          const nickname = sessionStorage.getItem('nickname');
          await this.$axios.put(`/conversations/${chat.id}/letti`, {}, {
          headers: { Authorization: "Bearer " + nickname }
            });
        } catch (error) {
           console.error("Errore update letto:", error);
           chat.lettura = null; //  rollback
          }
         }
        }
      this.showMess(this.chatAttiva); 
    },
      
    async getAllUsers() {
    try {
      // Sostituisci con la tua rotta reale per prendere tutti gli utenti
      const res = await this.$axios.get('/all-users'); 
      this.allUsers = res.data; 
    } catch (e) {
      console.error("Errore caricamento lista utenti", e);
    }
  }, 
  
    async showMess(convoId){
      try {
        sessionStorage.setItem("convoId", convoId);
        const nik = sessionStorage.getItem('nickname');
        const response = await this.$axios.get(`/conversations/${convoId}`, {
          headers: { Authorization: "Bearer " + nik }
            }); 
        const temp = Array.isArray(response.data) ? response.data : [];     /* questo per visualizzare chat vuote */
     
         const messConNick = await Promise.all(
         temp.map(async (msg) => {
              msg.sender_name = await this.getNickById(msg.sender);
              msg.forwarded = msg.forwarded === 1 || msg.forwarded === true;
              return msg;
              })
          );
        const messMap = {};

        // codice per rendere visibile la risposta
        messConNick.forEach(m => messMap[m.id] = m);
        messConNick.forEach(m => {
            // Se è una risposta (reply != 0 e reply != null)
            if (m.reply && messMap[m.reply]) {
                m.replyText = messMap[m.reply].text || "Foto/Media";
                m.replyUser = messMap[m.reply].sender_name;
            } else if (m.reply) {
                m.replyText = "Messaggio non trovato";
            }
        });
        this.messages = messConNick; /* così ho anche il nickname associato */
        this.numMess = messConNick.length;

        await Promise.all(messConNick.map(msg => this.getReacts(msg.id))); 

      }catch (error) {
        this.errormsg = "Errore nel recupero dei messaggi della chat: " + error.message;
        }
      },

    async inviaMess(convoId, testo, file, isRisp){  // ovvero lo aggiungo al database
    // in input vorrà convo ID etc
          try {
      const payload = {
        convoid: convoId,
        testo: testo || "",
        content: null,
        media_type: "text",
        file_type: null,
        reply: isRisp ? isRisp.id : 0
      };

    if (file) {
      // converto file in base64
      const base64 = await new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = e => resolve(e.target.result.split(',')[1]);
        reader.onerror = e => reject(e);
        reader.readAsDataURL(file);
      });

      payload.content = base64;
      payload.media_type = "image";
      payload.file_type = file.type;
    }
          const nik = sessionStorage.getItem('nickname');
          const response = await this.$axios.post(`/conversations/${convoId}/messages`,payload, {
            headers: { Authorization: "Bearer " + nik }
            }); 
          const messId = response.data.idMess;

          await this.showMess(convoId);
          this.isRisp = null;
        }catch (error) {
          this.errormsg = "Errore nel recupero del messaggio aggiunto: " + error.message;
        }
    },

    async cancelMess(messageId) {  // scritto qui per non dare problemi a prop
      if (!confirm("Sei sicuro di voler eliminare questo messaggio?")) return;

      
      try {
        const nickname = sessionStorage.getItem('nickname');
        const convoId = sessionStorage.getItem('convoId');
        await this.$axios.delete(`/conversations/${convoId}/messages/${messageId}`, {
          headers: { Authorization: "Bearer " + nickname }
        });
         const msgIndex = this.messages.findIndex(m => m.id === messageId);
          if (msgIndex !== -1) {
            this.messages[msgIndex] = {
              ...this.messages[msgIndex],
              text: null,
              deleted: true
            };
          }

        this.idMess = null;
      } catch (err) {
        alert("Errore nell'eliminazione del messaggio: " + err.message);
      }
    },
 // questo per aggiornare l'avatar di un gruppo
     updateAvatar ({ convId, avatar }){
      if (!avatar) return;
      
      const reader = new FileReader(); // conversione in base 64
      reader.onload= async (e) => {
        try {
          const immagineCompleta = e.target.result;
          const base64Only = immagineCompleta.split(',')[1];
          const payload = { photo: base64Only };

          const response = await this.$axios.put(`/groups/${convId}/photo`, payload);

          const index = this.chats.findIndex(c => c.id == convId);
      if (index !== -1) {
          // Nota: qui probabilmente stavi aggiornando 'avatar' correttamente
          this.chats[index] = { ...this.chats[index], avatar: immagineCompleta };
          this.chats = [...this.chats]; // Refresh array
      }

     // l'aggiornamento dell'interfaccia avviene prima
        
        } catch (err) {
          alert("Errore nell'aggiornamento della propic: " + err.message);
        }
      };
      reader.onerror = (err) => console.error('FileReader errore:', err);
      reader.readAsDataURL(avatar);
  
    },
// per cambiare il nome di un gruppo
    async updateName({convId,nome}){   // nome inserito sarà poi il payload
      try {
        const payload = { Nome: nome };
        const response = await this.$axios.put(`/groups/${convId}/name`,payload)

        const group = this.chats.find(g => g.id === convId);
        if (group) group.nome = nome;    
        }catch (error) {
          this.errormsg = "Errore nell'aggiornamento del nome : " + error.response;
        }

    },

    async addMember({convId,user}){
      try {
        const nik = sessionStorage.getItem('nickname');
        const payload = { User: user };
        const response = await this.$axios.post(`/groups/${convId}/members`,payload, {
            headers: { Authorization: "Bearer " + nik }
            });    
          alert("Membro aggiunto!");
        }catch (error) {
          if (error.response) {
            if (error.response.status === 409) {
                alert("Attenzione: " + error.response.data.error); 
                
            } else if (error.response.status === 404) {
                // utente inesistente
                alert("Errore: l'utente non esiste.");

            } else {
                // 500
                console.error("Errore server:", error.response);
                alert("Si è verificato un errore imprevisto: " + error.response.status);
            }
          } 
        }
    },

    async postReact(emoji, messageId) {
      const payload = {
          idEmoji: emoji,    /* ricorda che questo è il nome da json */
      };
      try {
          const nik = sessionStorage.getItem('nickname');
          const response = await this.$axios.post(`/messages/${messageId}/comments`,payload, {
            headers: { Authorization: "Bearer " + nik }
            });    
            await this.getReacts(messageId);
        }catch (error) {
          alert("Non puoi mettere due emoji!")
          this.errormsg = "Errore nella creazione delle emoji : " + error.message;          
        }
    },
    async deleteReact(emoji, messageId)  {
      
      try {
          const nik = sessionStorage.getItem('nickname');
          const response = await this.$axios.delete(`/messages/${messageId}/comments/${emoji}`, {
            headers: { Authorization: "Bearer " + nik }
            });    
             this.reactions[messageId] = this.reactions[messageId].filter(r => r.emoji !== emoji || r.user !== nik);
        }catch (error) {
          this.errormsg = "Errore nella cancellazione dell'emoji : " + error.message;
        }
    },

    checkReact(emoji, messageId) { /* le reazioni sono un array con due indici */
      const commentato = this.reactions[messageId]?.some(r => r.emoji === emoji);

      if (!commentato) {
        return this.postReact(emoji, messageId);
      } else {
        return this.deleteReact(emoji, messageId);
      }
    },
    async getReacts(messageId) {
      try {
          const nik = sessionStorage.getItem('nickname');
          const response = await this.$axios.get(`/messages/${messageId}/comments`);    
          const mapped = response.data.map(react => ({
            emoji: react.emoji,
            user: react.user
          }));
          this.reactions[String(messageId)] = mapped;
        }catch (error) {
          this.errormsg = "Errore nel recupero delle emoji : " + error.message;
          return [];
        }
    },
    inoltro({ convoId, newMessage }) {
       if (this.chatAttiva === convoId) {
        this.messages.push(newMessage);
      }
    },
},
  mounted() {
    this.showConvos();
    this.getAllUsers();
    this.intervalId = setInterval(() => {
    if (this.chatAttiva) {
      this.showMess(this.chatAttiva);
    } this.showConvos();  // aggiornamento delle conversazioni ogni 10 secondi  
      }, 10000);
    },
    beforeUnmount() {
  clearInterval(this.intervalId);
}
};
</script>

<template>
  <div class="d-flex">
     <SideConvo
      :visible="showConv"
      :chats="chats"
      @select-chat="selezionaChat"
    /> 
    
<!-- per aggiornare vista chat -->
    <template v-if="interaChat">
        <ConvoInterface
          :visible="true"
          :messages="messages"
          :myId="userId"
          :info="interaChat"  
          :convos="chats"
          :reactions="reactions"
          :allUsers="allUsers"
          @send-combined="inviaMess"
          @delete-mess="cancelMess"
          @update-avatar="updateAvatar"
          @rename-chat="updateName"
          @new-member="addMember"
          @check-react="({ emoji, messageId }) => checkReact(emoji, messageId)"
          @get-emoji="({ messageId }) => getReacts(messageId)"
          @message-forwarded="inoltro"/>
        <!-- importante per allineare l'id -->
      </template>

      <template v-else>
        <div class="main-area"> <!-- definisco una formattazione senza messaggi-->
           <h3>Seleziona una conversazione per iniziare</h3>
        </div>
      </template>
  </div>
</template>



<style scoped>
.conversation-container {
  height: 100vh;
}

.sidebar {
  width: 300px;
  border-right: 1px solid #ccc;
  overflow-y: auto;
  display: flex;

}

.main-area {
  flex: 1;
  color: orange;
  font-style: italic;
  padding: 20rem;
  padding-left: 40rem;
}

.placeholder {
  text-align: center;
  color: #888;
  margin-top: 50px;
}

</style>
