<script lang="ts">
	import { onMount } from 'svelte';
	import { createGuest } from '$lib/api';

	// Form state
	let formName = $state('');
	let formEmail = $state('');
	let formPhone = $state('');
	let formStatus = $state<'attending' | 'pending' | 'declined'>('attending');
	
	// UI state
	let formError = $state<string | null>(null);
	let formSubmitting = $state(false);
	let showSuccess = $state(false);
	let fieldErrors = $state<Record<string, string>>({});
let stats = $state({ total: 0, attending: 0, pending: 0, declined: 0 });

	onMount(async () => {
		await refreshStats();
	});

	async function refreshStats() {
		try {
			const response = await fetch('http://localhost:8080/api/guests/stats');
			if (response.ok) {
				const data = await response.json();
				// Update stats reactively
				stats.total = data.total || 0;
				stats.attending = data.attending || 0;
				stats.pending = data.pending || 0;
				stats.declined = data.declined || 0;
				
				console.log('Stats updated:', stats);
			}
		} catch (error) {
			console.error('Failed to fetch stats:', error);
		}
	}

	// Validation functions
	function validateEmail(email: string): boolean {
		const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
		return emailRegex.test(email);
	}

	function validatePhone(phone: string): boolean {
		if (!phone) return true; // Phone is optional
		const phoneRegex = /^[\d\s\-\+\(\)]{10,}$/;
		return phoneRegex.test(phone);
	}

	function validateForm(): boolean {
		const errors: Record<string, string> = {};

		// Name validation
		if (!formName.trim()) {
			errors.name = 'Please enter your name';
		} else if (formName.trim().length < 2) {
			errors.name = 'Name must be at least 2 characters';
		}

		// Email validation
		if (!formEmail.trim()) {
			errors.email = 'Please enter your email address';
		} else if (!validateEmail(formEmail)) {
			errors.email = 'Please enter a valid email address';
		}

		// Phone validation (optional but must be valid if provided)
		if (formPhone && !validatePhone(formPhone)) {
			errors.phone = 'Please enter a valid phone number';
		}

		fieldErrors = errors;
		return Object.keys(errors).length === 0;
	}

	// Real-time validation on blur
	function handleBlur(field: string) {
		const errors = { ...fieldErrors };
		
		switch (field) {
			case 'name':
				if (!formName.trim()) {
					errors.name = 'Please enter your name';
				} else if (formName.trim().length < 2) {
					errors.name = 'Name must be at least 2 characters';
				} else {
					delete errors.name;
				}
				break;
			case 'email':
				if (!formEmail.trim()) {
					errors.email = 'Please enter your email address';
				} else if (!validateEmail(formEmail)) {
					errors.email = 'Please enter a valid email address';
				} else {
					delete errors.email;
				}
				break;
			case 'phone':
				if (formPhone && !validatePhone(formPhone)) {
					errors.phone = 'Please enter a valid phone number';
				} else {
					delete errors.phone;
				}
				break;
		}
		
		fieldErrors = errors;
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		
		// Reset messages
		formError = null;
		showSuccess = false;

		// Validate form
		if (!validateForm()) {
			formError = 'Please fix the errors below';
			return;
		}

		try {
			formSubmitting = true;
			
			await createGuest({
				name: formName.trim(),
				email: formEmail.trim(),
				phone: formPhone.trim() || '',
				status: formStatus,
							

			});

			// Show success message
			showSuccess = true;
		await refreshStats(); 
			// Reset form
			formName = '';
			formEmail = '';
			formPhone = '';
			formStatus = 'attending';
			fieldErrors = {};
			
			// Scroll to success message
			window.scrollTo({ top: 0, behavior: 'smooth' });
			
		} catch (e) {
			formError = e instanceof Error ? e.message : 'Unable to submit your RSVP. Please try again.';
			window.scrollTo({ top: 0, behavior: 'smooth' });
		} finally {
			formSubmitting = false;
		}
	}

	function getStatusLabel(status: string): string {
		switch (status) {
			case 'attending':
				return "Yes, I'll be there! 🎉";
			case 'pending':
				return "Maybe 🤔";
			case 'declined':
				return "Sorry, can't make it 😔";
			default:
				return status;
		}
	}

	function getStatusDescription(status: string): string {
		switch (status) {
			case 'attending':
				return "We can't wait to see you!";
			case 'pending':
				return "Let us know when you decide";
			case 'declined':
				return "We'll miss you!";
			default:
				return '';
		}
	}
</script>

<svelte:head>
	<title>RSVP - Event Invitation</title>
	<meta name="description" content="RSVP for our upcoming event" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-purple-50 via-pink-50 to-blue-50">
	<!-- Hero Section -->
	<div class="bg-gradient-to-r from-purple-600 to-pink-600 text-white">
		<div class="max-w-4xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
		
			<div class="text-center">
				<h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold mb-4 animate-fade-in">
					You're Invited!
				</h1>
				<p class="text-xl sm:text-2xl text-purple-100 mb-2">
					Join us for an unforgettable celebration
				</p>
				<div class="mt-8 space-y-2 text-lg">
					<div class="flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clip-rule="evenodd" />
						</svg>
						<span>Saturday, December 15th, 2025 at 6:00 PM</span>
					</div>
					<div class="flex items-center justify-center gap-2">
						<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
						</svg>
						<span>Grand Ballroom, Downtown Convention Center</span>
					</div>
				</div>
				
			</div>
	
		</div>
<div class="bg-white/20 backdrop-blur-sm text-white px-6 py-4 rounded-lg text-center mx-auto max-w-md mt-8">
    <h3 class="text-2xl font-bold">🎉 {stats.attending} people attending!</h3>
    <p class="text-sm mt-1 text-purple-100">Join them at our amazing event!</p>
</div>


	</div>

	<!-- RSVP Form Section -->
	<div class="max-w-2xl mx-auto px-4 py-12 sm:px-6 lg:px-8">
		<!-- Success Message -->
		{#if showSuccess}
			<div class="bg-green-50 border-2 border-green-400 rounded-lg p-6 mb-8 animate-slide-down">
				<div class="flex items-start">
					<div class="flex-shrink-0">
						<svg class="h-8 w-8 text-green-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-4">
						<h3 class="text-lg font-semibold text-green-800 mb-1">
							RSVP Submitted Successfully! 🎉
						</h3>
						<p class="text-green-700">
							Thank you for your response! We've sent a confirmation to your email address.
						</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- Error Message -->
		{#if formError}
			<div class="bg-red-50 border-2 border-red-400 rounded-lg p-4 mb-8 animate-slide-down">
				<div class="flex items-start">
					<div class="flex-shrink-0">
						<svg class="h-6 w-6 text-red-400" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<p class="text-sm font-medium text-red-800">{formError}</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- RSVP Form Card -->
		<div class="bg-white rounded-2xl shadow-xl overflow-hidden">
			<div class="bg-gradient-to-r from-purple-500 to-pink-500 px-6 py-4">
				<h2 class="text-2xl font-bold text-white">Please RSVP</h2>
				<p class="text-purple-100 text-sm mt-1">Let us know if you can join us</p>
			</div>

			<form onsubmit={handleSubmit} class="p-6 sm:p-8 space-y-6">
				<!-- Name Field -->
				<div>
					<label for="name" class="block text-sm font-semibold text-gray-700 mb-2">
						Your Name <span class="text-red-500">*</span>
					</label>
					<input
						type="text"
						id="name"
						bind:value={formName}
						onblur={() => handleBlur('name')}
						class="w-full px-4 py-3 border-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all {fieldErrors.name ? 'border-red-300 bg-red-50' : 'border-gray-300'}"
						placeholder="John Doe"
						disabled={formSubmitting}
					/>
					{#if fieldErrors.name}
						<p class="mt-2 text-sm text-red-600 flex items-center gap-1">
							<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
							</svg>
							{fieldErrors.name}
						</p>
					{/if}
				</div>

				<!-- Email Field -->
				<div>
					<label for="email" class="block text-sm font-semibold text-gray-700 mb-2">
						Email Address <span class="text-red-500">*</span>
					</label>
					<input
						type="email"
						id="email"
						bind:value={formEmail}
						onblur={() => handleBlur('email')}
						class="w-full px-4 py-3 border-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all {fieldErrors.email ? 'border-red-300 bg-red-50' : 'border-gray-300'}"
						placeholder="john@example.com"
						disabled={formSubmitting}
					/>
					{#if fieldErrors.email}
						<p class="mt-2 text-sm text-red-600 flex items-center gap-1">
							<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
							</svg>
							{fieldErrors.email}
						</p>
					{/if}
				</div>

				<!-- Phone Field -->
				<div>
					<label for="phone" class="block text-sm font-semibold text-gray-700 mb-2">
						Phone Number <span class="text-gray-400 text-xs">(Optional)</span>
					</label>
					<input
						type="tel"
						id="phone"
						bind:value={formPhone}
						onblur={() => handleBlur('phone')}
						class="w-full px-4 py-3 border-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all {fieldErrors.phone ? 'border-red-300 bg-red-50' : 'border-gray-300'}"
						placeholder="+1 (555) 123-4567"
						disabled={formSubmitting}
					/>
					{#if fieldErrors.phone}
						<p class="mt-2 text-sm text-red-600 flex items-center gap-1">
							<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
							</svg>
							{fieldErrors.phone}
						</p>
					{/if}
				</div>

				<!-- RSVP Status Selection -->
				<div>
					<label class="block text-sm font-semibold text-gray-700 mb-3">
						Will you be attending? <span class="text-red-500">*</span>
					</label>
					<div class="space-y-3">
						{#each ['attending', 'pending', 'declined'] as status}
							<label
								class="relative flex items-start p-4 border-2 rounded-lg cursor-pointer transition-all hover:bg-gray-50 {formStatus === status ? 'border-purple-500 bg-purple-50' : 'border-gray-200'}"
							>
								<input
									type="radio"
									name="status"
									value={status}
									bind:group={formStatus}
									class="mt-0.5 h-5 w-5 text-purple-600 focus:ring-purple-500 cursor-pointer"
									disabled={formSubmitting}
								/>
								<div class="ml-3 flex-1">
									<span class="block text-base font-medium text-gray-900">
										{getStatusLabel(status)}
									</span>
									<span class="block text-sm text-gray-500 mt-0.5">
										{getStatusDescription(status)}
									</span>
								</div>
							</label>
						{/each}
					</div>
				</div>

				<!-- Submit Button -->
				<div class="pt-4">
					<button
						type="submit"
						disabled={formSubmitting}
						class="w-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 disabled:from-gray-400 disabled:to-gray-400 text-white font-semibold py-4 px-6 rounded-lg shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200 disabled:cursor-not-allowed disabled:transform-none flex items-center justify-center gap-2"
					>
						{#if formSubmitting}
							<svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							<span>Submitting...</span>
						{:else}
							<span>Submit RSVP</span>
							<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
							</svg>
						{/if}
					</button>
				</div>

				<p class="text-center text-sm text-gray-500 pt-2">
					<span class="text-red-500">*</span> Required fields
				</p>
			</form>
		</div>

		<!-- Additional Info Section -->
		<div class="mt-12 bg-white rounded-xl shadow-md p-6 sm:p-8">
			<h3 class="text-xl font-bold text-gray-900 mb-4 flex items-center gap-2">
				<svg class="w-6 h-6 text-purple-600" fill="currentColor" viewBox="0 0 20 20">
					<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
				</svg>
				Event Details
			</h3>
			<div class="space-y-4 text-gray-600">
				<div class="flex items-start gap-3">
					<svg class="w-5 h-5 text-purple-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
					</svg>
					<div>
						<p class="font-semibold text-gray-900">Time</p>
						<p>Doors open at 5:30 PM • Event starts at 6:00 PM</p>
					</div>
				</div>
				<div class="flex items-start gap-3">
					<svg class="w-5 h-5 text-purple-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
					</svg>
					<div>
						<p class="font-semibold text-gray-900">Location</p>
						<p>Grand Ballroom, Downtown Convention Center</p>
						<p class="text-sm text-gray-500">123 Main Street, City, State 12345</p>
					</div>
				</div>
				<div class="flex items-start gap-3">
					<svg class="w-5 h-5 text-purple-600 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
						<path d="M2 3a1 1 0 011-1h2.153a1 1 0 01.986.836l.74 4.435a1 1 0 01-.54 1.06l-1.548.773a11.037 11.037 0 006.105 6.105l.774-1.548a1 1 0 011.059-.54l4.435.74a1 1 0 01.836.986V17a1 1 0 01-1 1h-2C7.82 18 2 12.18 2 5V3z" />
					</svg>
					<div>
						<p class="font-semibold text-gray-900">Questions?</p>
						<p>Contact us at events@example.com or (555) 123-4567</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	@keyframes fade-in {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes slide-down {
		from {
			opacity: 0;
			transform: translateY(-20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.animate-fade-in {
		animation: fade-in 0.6s ease-out;
	}

	.animate-slide-down {
		animation: slide-down 0.4s ease-out;
	}
</style>
