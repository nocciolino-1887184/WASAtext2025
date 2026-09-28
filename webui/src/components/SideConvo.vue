<script>
export default {
  name: "SidebarConversations",
  props: {
    visible: Boolean,
    chats: Array
  },
  computed: {   // computed è utile per fare cache dei dati
     formattedChats() {
       return this.chats.map(chat => {        
        return {
          id: chat.id,
          name: chat.name || chat.ricevitore || "Sconosciuto",
          avatar: chat.avatar,
          lastMessage: chat.lastMessage || "Ancora nessun messaggio",
          time:  chat.time || null,   // il tempo qui ha un formato diverdso
          lettura: chat.lettura,
          tipo: chat.tipo || null
        };
    }); // l'ora forse non si vede per il poco spazio
  }

  },
  methods: {
    formatTime(aString) {
      if (!aString) return '';
      
      const date = new Date(aString + 'Z'); /* questo serve per assicurarmi sia in UTC */
      if (isNaN(date.getTime())) return 'Data non valida';
      return date.toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit', timeZone: 'Europe/Rome' }); /* conversione corretta */      
    },

    async trasmetti(chat) {
        chat.lettura = null;   // se ci clicco diventa bianco
        this.$emit('select-chat', chat);
       }
   }
}
</script>

<template>
  <div class="sidebar" v-show="visible">
    <div class="sidebar-header">
      <h5>Conversazioni Esistenti</h5>
    </div>

    <div class="conversation" v-for="chat in formattedChats" :key="chat.id" :class="{ 'non-letto': chat.lettura }" @click="trasmetti(chat)">
      <img :src="chat.avatar || '/default-profile.jpg'"  @error="e => e.target.src = '/default-profile.jpg'" class="avatar" alt="avatar" />
  
      <div class="info">
        <div class="name">{{ chat.name }}</div>
        <div class="last-message">  
          <span class="text">{{ chat.lastMessage }}</span> 
          <span class="time">{{ formatTime(chat.time) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>


<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 320px; /*attaccato all'altra barra */
  border-left: 1px solid #ccc;
  width: 350px;
  height: 100vh;
  background-color: #f8f9fa;
  box-shadow: 2px 0 5px rgba(0,0,0,0.1);
  padding: 1rem;
  overflow-y: auto;
  z-index: 999;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding-top: 3rem ;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

.conversation {
  display: flex;
  align-items: center;
  margin-top: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 1rem;
}

.info {
  flex-grow: 1;
}

.name {
  font-weight: bold;
  width: 180px;  /*scrivere questo permette di tagliare*/
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;  
}

.last-message {
  display: flex; 
  position: relative;    
  font-size: 0.85rem;
  color: #666;          
  max-width: 200px;    
}

.time {
  font-size: 0.75rem;
  bottom: 0;
  right: 0;
}

.last-message .text {
  flex: 1;
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  padding-bottom: 1.2rem;
  min-width: 0;
}

.last-message .time {
  flex-shrink: 0;
  position: absolute; 
  margin-left: 0.5rem;
  font-size: 0.75rem;
}

.conversation.non-letto {
  background-color: #fafbb6; 
}


</style>
