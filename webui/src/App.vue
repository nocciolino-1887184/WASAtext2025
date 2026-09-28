<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	methods: {
		logout() {			// il logout qui è accessibile da tutti
			const conferma = window.confirm("Sei sicuro di voler uscire?");
			if (conferma) {
			// chiusura e cancellazione della sessione
				sessionStorage.clear();
				}
			alert("Logout avvenuto con successo.");

			this.$router.push('/session');
		}
	} 
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/dashboard">WASA Text</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/dashboard" class="nav-link" active-class="active">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Homepage
							</RouterLink>
						</li>
						<!-- <li class="nav-item">
							<RouterLink to="/session" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Login
							</RouterLink>
						</li> -->
						<li class="nav-item">
							<RouterLink to="/users/:id" class="nav-link" active-class="active">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
								Profilo
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users/:id/conversations" class="nav-link" active-class="active">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
								Conversazioni
							</RouterLink>
						</li>
						<li class="nav-item">
							<a href="#" class="nav-link" @click.prevent="logout">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#power"/></svg>
								Logout
							</a>
						</li>
					</ul>

				<!-- <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Options</span>
					</h6> 
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link" active-class="active">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
								Settings
							</RouterLink> 
						</li>
					</ul>-->
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4"> 
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>
