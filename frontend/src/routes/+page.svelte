<script>
	import { onMount } from 'svelte';
	
	let searchQuery = '';
	let results = [];
	let isLoading = false;
	let error = '';
	let showToast = false;
	let toastMessage = '';
	let mode = 'browse'; // 'search' | 'browse'
	
	// Refinements state
	let selectedCategories = [];
	let selectedTags = [];
	let sort = 'popularity';
	let order = 'desc';
	let page = 1;
	let limit = 20;
	let showRefinements = false;
	
	// Facets data
	let facets = {};
	let pagination = {};
	let mounted = false;

	async function loadQuotes() {
		if (!mounted) return;
		
		// Determine mode based on search query
		mode = searchQuery.trim() ? 'search' : 'browse';
		isLoading = true;
		error = '';
		
		try {
			const params = new URLSearchParams({
				page: page.toString(),
				limit: limit.toString(),
				sort,
				order,
				facets: 'true'
			});
			
			// Add search query if present
			if (searchQuery.trim()) {
				params.set('q', searchQuery.trim());
			}
			
			// Add array filters
			selectedCategories.forEach(cat => params.append('categories[]', cat));
			selectedTags.forEach(tag => params.append('tags[]', tag));
			
			// Use unified search API for both search and browse
			const response = await fetch(`http://localhost:8080/api/search?${params}`);
			
			if (!response.ok) {
				throw new Error('Request failed');
			}
			
			const data = await response.json();
			results = data.quotes || data.results || [];
			facets = data.facets || {};
			pagination = data.pagination || { total_count: data.count || 0 };
		} catch (err) {
			error = err.message || 'An error occurred';
			results = [];
		} finally {
			isLoading = false;
		}
	}

	// Simplified search handler - just calls loadQuotes
	async function handleSearch() {
		page = 1; // Reset to first page when searching
		return await loadQuotes();
	}

	function copyToClipboard(text) {
		navigator.clipboard.writeText(text).then(() => {
			toastMessage = 'Quote copied to clipboard!';
			showToast = true;
			setTimeout(() => showToast = false, 3000);
		});
	}

	function toggleCategory(category) {
		if (selectedCategories.includes(category)) {
			selectedCategories = selectedCategories.filter(c => c !== category);
		} else {
			selectedCategories = [...selectedCategories, category];
		}
		page = 1;
		searchQuery = ''; // Clear search when using filters
		loadQuotes();
	}

	function toggleTag(tag) {
		if (selectedTags.includes(tag)) {
			selectedTags = selectedTags.filter(t => t !== tag);
		} else {
			selectedTags = [...selectedTags, tag];
		}
		page = 1;
		searchQuery = ''; // Clear search when using filters
		loadQuotes();
	}

	function clearAllFilters() {
		selectedCategories = [];
		selectedTags = [];
		sort = 'popularity';
		order = 'desc';
		page = 1;
		loadQuotes();
	}

	function changePage(newPage) {
		page = newPage;
		loadQuotes();
	}

	// Load initial quotes on page load
	onMount(() => {
		mounted = true;
		loadQuotes();
	});
</script>

<!-- App Shell -->
<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100">
	<!-- Header -->
	<header class="relative overflow-hidden">
		<!-- Background Pattern -->
		<div class="absolute inset-0 bg-gradient-to-r from-blue-600 via-indigo-600 to-purple-600"></div>
		<div class="absolute inset-0 bg-black/20"></div>
		<div class="absolute inset-0 opacity-30">
			<div class="w-full h-full" style="background-image: radial-gradient(circle at 25px 25px, rgba(255,255,255,0.1) 2px, transparent 0); background-size: 50px 50px;"></div>
		</div>
		
		<div class="relative container mx-auto px-4 py-12">
			<div class="text-center text-white">
				<div class="inline-flex items-center justify-center w-20 h-20 bg-white/20 backdrop-blur-sm rounded-full mb-6">
					<span class="text-4xl">📚</span>
				</div>
				<h1 class="text-5xl md:text-6xl font-bold mb-4 leading-tight">
					<span class="bg-gradient-to-r from-white to-blue-100 bg-clip-text text-transparent">
						Themed Quote Collections
					</span>
				</h1>
				<p class="text-xl text-blue-100 max-w-2xl mx-auto leading-relaxed">
					Discover inspiring quotes by themes, authors, and keywords
				</p>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<main class="container mx-auto px-4 -mt-8 relative z-10 space-y-8">
		<!-- Search Section -->
		<div class="bg-white/80 backdrop-blur-lg rounded-2xl shadow-xl border border-white/20 p-8">
			<form on:submit|preventDefault={handleSearch} class="space-y-6">
				<div class="space-y-3">
					<label class="block text-lg font-semibold text-gray-800" for="search">
						🔍 Search Quotes
					</label>
					<div class="relative">
						<div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
							<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
							</svg>
						</div>
						<input
							id="search"
							bind:value={searchQuery}
							type="text"
							placeholder="Enter themes like 'courage', 'innovation', or author names..."
							class="w-full pl-12 pr-4 py-4 text-lg border-2 border-gray-200 rounded-xl focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 bg-white/90"
						/>
						<button 
							type="submit" 
							disabled={isLoading}
							class="absolute right-2 top-2 px-6 py-2 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 font-medium shadow-lg hover:shadow-xl transform hover:scale-105"
						>
							{isLoading ? '⏳ Searching...' : '✨ Search'}
						</button>
					</div>
				</div>
			</form>
		</div>

		<!-- Refinements Section -->
		<div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-white/20 p-6">
			<div class="flex items-center justify-between mb-6">
				<h3 class="text-xl font-bold text-gray-800 flex items-center gap-2">
					🎯 Refinements
				</h3>
				<button 
					on:click={() => showRefinements = !showRefinements}
					class="px-4 py-2 text-sm bg-gradient-to-r from-gray-100 to-gray-200 hover:from-gray-200 hover:to-gray-300 text-gray-700 rounded-lg transition-all duration-200 font-medium shadow-sm hover:shadow-md transform hover:scale-105"
				>
					{showRefinements ? '🔺 Hide' : '🔻 Show'} Filters
				</button>
			</div>

			{#if showRefinements}
				<div class="space-y-8">
					<!-- Sorting Controls -->
					<div class="flex flex-wrap gap-6">
						<div class="flex items-center gap-3">
							<label for="sort-select" class="text-sm font-semibold text-gray-700">📊 Sort by:</label>
							<select 
								id="sort-select"
								bind:value={sort} 
								on:change={() => { page = 1; loadQuotes(); }}
								class="px-4 py-2 border-2 border-gray-200 rounded-lg text-sm focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 bg-white/90 font-medium"
							>
								<option value="popularity">⭐ Popularity</option>
								<option value="created_at">📅 Date Added</option>
								<option value="random">🎲 Random</option>
							</select>
						</div>
						<div class="flex items-center gap-3">
							<label for="order-select" class="text-sm font-semibold text-gray-700">🔄 Order:</label>
							<select 
								id="order-select"
								bind:value={order} 
								on:change={() => { page = 1; loadQuotes(); }}
								class="px-4 py-2 border-2 border-gray-200 rounded-lg text-sm focus:ring-4 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 bg-white/90 font-medium"
							>
								<option value="desc">📈 High to Low</option>
								<option value="asc">📉 Low to High</option>
							</select>
						</div>
						{#if selectedCategories.length > 0 || selectedTags.length > 0}
							<button 
								on:click={clearAllFilters}
								class="px-4 py-2 text-sm bg-gradient-to-r from-red-500 to-pink-500 hover:from-red-600 hover:to-pink-600 text-white rounded-lg transition-all duration-200 font-medium shadow-lg hover:shadow-xl transform hover:scale-105"
							>
								🗑️ Clear All
							</button>
						{/if}
					</div>

					<!-- Active Filters -->
					{#if selectedCategories.length > 0 || selectedTags.length > 0}
						<div class="space-y-3">
							<h4 class="text-sm font-semibold text-gray-700 flex items-center gap-2">✅ Active Filters:</h4>
							<div class="flex flex-wrap gap-3">
								{#each selectedCategories as category}
									<span class="inline-flex items-center gap-2 bg-gradient-to-r from-purple-500 to-indigo-500 text-white text-sm px-4 py-2 rounded-full shadow-md hover:shadow-lg transition-all duration-200">
										<span class="font-medium">📁 {category}</span>
										<button on:click={() => toggleCategory(category)} class="text-white/80 hover:text-white hover:bg-white/20 rounded-full p-1 transition-all duration-200">
											<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
										</button>
									</span>
								{/each}
								{#each selectedTags as tag}
									<span class="inline-flex items-center gap-2 bg-gradient-to-r from-blue-500 to-cyan-500 text-white text-sm px-4 py-2 rounded-full shadow-md hover:shadow-lg transition-all duration-200">
										<span class="font-medium">🏷️ {tag}</span>
										<button on:click={() => toggleTag(tag)} class="text-white/80 hover:text-white hover:bg-white/20 rounded-full p-1 transition-all duration-200">
											<svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path></svg>
										</button>
									</span>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Category Facets -->
					{#if facets.categories && facets.categories.length > 0}
						<div class="space-y-4">
							<h4 class="text-sm font-semibold text-gray-700 flex items-center gap-2">📂 Categories:</h4>
							<div class="flex flex-wrap gap-3">
								{#each facets.categories as category}
									<button 
										on:click={() => toggleCategory(category.value)}
										class="inline-flex items-center gap-2 px-4 py-2 text-sm rounded-xl transition-all duration-200 font-medium shadow-sm hover:shadow-md transform hover:scale-105 {selectedCategories.includes(category.value) ? 'bg-gradient-to-r from-purple-600 to-indigo-600 text-white shadow-lg' : 'bg-white/80 hover:bg-white text-gray-700 border-2 border-gray-200 hover:border-purple-300'}"
									>
										<span>{category.value}</span>
										<span class="text-xs bg-white/20 px-2 py-0.5 rounded-full {selectedCategories.includes(category.value) ? 'text-white/90' : 'text-gray-500 bg-gray-100'}">
											{category.count}
										</span>
									</button>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Tag Facets -->
					{#if facets.tags && facets.tags.length > 0}
						<div class="space-y-4">
							<h4 class="text-sm font-semibold text-gray-700 flex items-center gap-2">🏷️ Popular Tags:</h4>
							<div class="flex flex-wrap gap-3">
								{#each facets.tags as tag}
									<button 
										on:click={() => toggleTag(tag.value)}
										class="inline-flex items-center gap-2 px-4 py-2 text-sm rounded-xl transition-all duration-200 font-medium shadow-sm hover:shadow-md transform hover:scale-105 {selectedTags.includes(tag.value) ? 'bg-gradient-to-r from-blue-600 to-cyan-600 text-white shadow-lg' : 'bg-white/80 hover:bg-white text-gray-700 border-2 border-gray-200 hover:border-blue-300'}"
									>
										<span>{tag.value}</span>
										<span class="text-xs bg-white/20 px-2 py-0.5 rounded-full {selectedTags.includes(tag.value) ? 'text-white/90' : 'text-gray-500 bg-gray-100'}">
											{tag.count}
										</span>
									</button>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Loading Animation -->
		{#if isLoading}
			<div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-white/20 p-8">
				<div class="flex items-center justify-center space-x-4">
					<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
					<div class="text-lg font-medium text-gray-700">Searching for amazing quotes...</div>
				</div>
				<div class="mt-4 w-full bg-gradient-to-r from-blue-200 via-purple-200 to-indigo-200 rounded-full h-3">
					<div class="bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 h-3 rounded-full animate-pulse" style="width: 70%"></div>
				</div>
			</div>
		{/if}

		<!-- Error Message -->
		{#if error}
			<div class="bg-gradient-to-r from-red-50 to-pink-50 border-2 border-red-200 text-red-800 p-6 rounded-2xl shadow-lg">
				<div class="flex items-center space-x-3">
					<div class="flex-shrink-0">
						<div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
							<span class="text-xl">⚠️</span>
						</div>
					</div>
					<div>
						<h3 class="text-lg font-semibold">Oops! Something went wrong</h3>
						<p class="text-red-600 mt-1">{error}</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- Results -->
		{#if results.length > 0}
			<div class="space-y-8">
				<!-- Mamba UI Statistics Section -->
				<div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-white/20 p-8">
					<div class="text-center mb-8">
						<h2 class="text-3xl font-bold text-gray-800 mb-2">
							{mode === 'search' ? '🔍 Search Results' : '📖 Quote Collection'}
						</h2>
						<p class="text-gray-600">Discover wisdom from great minds throughout history</p>
					</div>
					
					<!-- Mamba UI Stats Grid -->
					<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
						<!-- Total Quotes -->
						<div class="bg-gradient-to-br from-violet-50 to-purple-50 rounded-2xl p-6 border border-violet-100">
							<div class="flex items-center justify-between">
								<div>
									<p class="text-sm font-medium text-violet-600 mb-1">Total Quotes</p>
									<p class="text-3xl font-bold text-gray-900">{pagination.total_count || results.length}</p>
								</div>
								<div class="w-12 h-12 bg-gradient-to-br from-violet-500 to-purple-600 rounded-xl flex items-center justify-center">
									<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
									</svg>
								</div>
							</div>
						</div>
						
						<!-- Current Page -->
						{#if mode === 'browse' && pagination.total_pages}
							<div class="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-2xl p-6 border border-blue-100">
								<div class="flex items-center justify-between">
									<div>
										<p class="text-sm font-medium text-blue-600 mb-1">Current Page</p>
										<p class="text-3xl font-bold text-gray-900">{pagination.page} of {pagination.total_pages}</p>
									</div>
									<div class="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center">
										<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
										</svg>
									</div>
								</div>
							</div>
						{:else}
							<div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-6 border border-green-100">
								<div class="flex items-center justify-between">
									<div>
										<p class="text-sm font-medium text-green-600 mb-1">Search Mode</p>
										<p class="text-xl font-bold text-gray-900">Active</p>
									</div>
									<div class="w-12 h-12 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center">
										<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
										</svg>
									</div>
								</div>
							</div>
						{/if}
						
						<!-- Categories/Tags Count -->
						<div class="bg-gradient-to-br from-orange-50 to-red-50 rounded-2xl p-6 border border-orange-100">
							<div class="flex items-center justify-between">
								<div>
									<p class="text-sm font-medium text-orange-600 mb-1">Categories</p>
									<p class="text-3xl font-bold text-gray-900">{facets.categories ? facets.categories.length : 0}</p>
								</div>
								<div class="w-12 h-12 bg-gradient-to-br from-orange-500 to-red-600 rounded-xl flex items-center justify-center">
									<svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
									</svg>
								</div>
							</div>
						</div>
					</div>
				</div>
				
				<!-- Mamba UI Enhanced Quote Cards -->
				<div class="grid gap-8">
					{#each results as quote}
						<!-- Mamba UI Testimonial Card Style -->
						<div class="group relative bg-white/90 backdrop-blur-sm rounded-3xl shadow-lg border border-white/20 hover:shadow-2xl transition-all duration-300 overflow-hidden">
							<!-- Gradient Background Accent -->
							<div class="absolute inset-0 bg-gradient-to-br from-violet-50 via-blue-50 to-indigo-50 opacity-30"></div>
							
							<!-- Quote Content -->
							<div class="relative p-8 lg:p-10">
								<!-- Quote Icon -->
								<div class="absolute top-6 left-6 w-12 h-12 bg-gradient-to-br from-violet-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg">
									<svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
										<path d="M14.017 21v-7.391c0-5.704 3.731-9.57 8.983-10.609l.995 2.151c-2.432.917-3.995 3.638-3.995 5.849h4v10h-9.983zm-14.017 0v-7.391c0-5.704 3.748-9.57 9-10.609l.996 2.151c-2.433.917-3.996 3.638-3.996 5.849h4v10h-10z"/>
									</svg>
								</div>
								
								<!-- Main Quote -->
								<div class="mt-8 mb-8">
									<blockquote class="text-xl lg:text-2xl font-semibold text-gray-800 leading-relaxed">
										{quote.quote}
									</blockquote>
								</div>
								
								<!-- Author Section -->
								<div class="flex items-center justify-between">
									<div class="flex items-center space-x-4">
										<!-- Author Avatar -->
										<div class="w-12 h-12 bg-gradient-to-br from-gray-200 to-gray-300 rounded-full flex items-center justify-center">
											<span class="text-lg font-bold text-gray-600">{quote.author.charAt(0)}</span>
										</div>
										<div>
											<p class="text-lg font-semibold text-gray-900">{quote.author}</p>
											{#if quote.category}
												<p class="text-sm text-gray-500 capitalize">{quote.category} Wisdom</p>
											{/if}
										</div>
									</div>
									
									<!-- Action Buttons -->
									<div class="flex items-center space-x-3">
										<!-- Share Button -->
										<button 
											on:click={() => copyToClipboard(`"${quote.quote}" - ${quote.author}`)}
											class="group/btn relative inline-flex items-center justify-center px-4 py-2 bg-white border border-gray-200 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-300 transition-all duration-200 shadow-sm"
										>
											<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
											</svg>
											Copy
										</button>
										
										<!-- Favorite Button -->
										<button class="group/fav relative inline-flex items-center justify-center p-2 bg-gradient-to-r from-violet-500 to-purple-600 hover:from-violet-600 hover:to-purple-700 rounded-xl text-white transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-105">
											<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
											</svg>
										</button>
									</div>
								</div>
								
								<!-- Tags Section -->
								{#if (quote.tags && quote.tags.length > 0) || quote.relevance || quote.popularity}
									<div class="mt-6 pt-6 border-t border-gray-100">
										<div class="flex flex-wrap gap-2">
											{#if quote.tags && quote.tags.length > 0}
												{#each quote.tags as tag}
													<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-violet-100 text-violet-800 hover:bg-violet-200 transition-colors cursor-pointer">
														{tag}
													</span>
												{/each}
											{/if}
											{#if quote.relevance}
												<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
													Score {quote.relevance.toFixed(2)}
												</span>
											{/if}
											{#if quote.popularity}
												<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
													⭐ {quote.popularity.toFixed(3)}
												</span>
											{/if}
										</div>
									</div>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<!-- Pagination -->
				{#if mode === 'browse' && pagination.total_pages > 1}
					<div class="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-white/20 p-6">
						<div class="flex justify-center items-center gap-4">
							<button 
								on:click={() => changePage(pagination.page - 1)}
								disabled={!pagination.has_prev}
								class="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-gray-100 to-gray-200 hover:from-blue-500 hover:to-indigo-500 text-gray-700 hover:text-white rounded-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:from-gray-100 disabled:hover:to-gray-200 disabled:hover:text-gray-700 transition-all duration-200 font-medium shadow-sm hover:shadow-lg transform hover:scale-105"
							>
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
								</svg>
								Previous
							</button>
							
							<div class="flex items-center gap-3">
								<span class="px-6 py-3 bg-gradient-to-r from-blue-500 to-indigo-500 text-white rounded-xl font-semibold shadow-lg">
									Page {pagination.page} of {pagination.total_pages}
								</span>
							</div>
							
							<button 
								on:click={() => changePage(pagination.page + 1)}
								disabled={!pagination.has_next}
								class="flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-gray-100 to-gray-200 hover:from-blue-500 hover:to-indigo-500 text-gray-700 hover:text-white rounded-xl disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:from-gray-100 disabled:hover:to-gray-200 disabled:hover:text-gray-700 transition-all duration-200 font-medium shadow-sm hover:shadow-lg transform hover:scale-105"
							>
								Next
								<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
								</svg>
							</button>
						</div>
					</div>
				{/if}
			</div>
		{:else if !isLoading && searchQuery && mode === 'search'}
			<div class="bg-white/90 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-12 text-center space-y-6">
				<div class="w-24 h-24 mx-auto bg-gradient-to-br from-gray-100 to-gray-200 rounded-full flex items-center justify-center">
					<span class="text-4xl">🔍</span>
				</div>
				<div class="space-y-3">
					<h3 class="text-2xl font-bold text-gray-800">No quotes found</h3>
					<p class="text-gray-600 text-lg max-w-md mx-auto leading-relaxed">
						No quotes found for "<strong class="text-blue-600">{searchQuery}</strong>". 
						Try searching for themes like "courage", "love", or author names.
					</p>
				</div>
				<div class="flex flex-wrap justify-center gap-3 pt-4">
					<span class="px-4 py-2 bg-blue-100 text-blue-800 rounded-full text-sm font-medium">💡 Try: "wisdom"</span>
					<span class="px-4 py-2 bg-purple-100 text-purple-800 rounded-full text-sm font-medium">💡 Try: "Einstein"</span>
					<span class="px-4 py-2 bg-green-100 text-green-800 rounded-full text-sm font-medium">💡 Try: "success"</span>
				</div>
			</div>
		{:else if !isLoading && mode === 'browse' && results.length === 0}
			<div class="bg-white/90 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-12 text-center space-y-6">
				<div class="w-24 h-24 mx-auto bg-gradient-to-br from-orange-100 to-red-100 rounded-full flex items-center justify-center">
					<span class="text-4xl">📂</span>
				</div>
				<div class="space-y-3">
					<h3 class="text-2xl font-bold text-gray-800">No quotes match your filters</h3>
					<p class="text-gray-600 text-lg max-w-md mx-auto leading-relaxed">
						Try adjusting your category or tag filters to see more results.
					</p>
				</div>
				<button 
					on:click={clearAllFilters}
					class="px-6 py-3 bg-gradient-to-r from-blue-500 to-indigo-500 hover:from-blue-600 hover:to-indigo-600 text-white rounded-xl font-semibold shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200"
				>
					🔄 Clear All Filters
				</button>
			</div>
		{/if}
	</main>

	<!-- Footer -->
	<footer class="mt-20 relative overflow-hidden">
		<div class="absolute inset-0 bg-gradient-to-r from-gray-900 via-blue-900 to-indigo-900"></div>
		<div class="absolute inset-0 bg-black/20"></div>
		<div class="relative py-12 text-center">
			<div class="container mx-auto px-4">
				<div class="inline-flex items-center justify-center w-16 h-16 bg-white/10 backdrop-blur-sm rounded-full mb-6">
					<span class="text-2xl">⚡</span>
				</div>
				<p class="text-white/80 text-lg max-w-2xl mx-auto leading-relaxed">
					Powered by <span class="font-semibold text-blue-300">ParadeDB</span> full-text search • 
					Built with <span class="font-semibold text-green-300">SvelteKit</span> & 
					<span class="font-semibold text-cyan-300">Tailwind CSS</span>
				</p>
				<div class="mt-6 flex justify-center space-x-6">
					<div class="w-2 h-2 bg-blue-400 rounded-full animate-pulse"></div>
					<div class="w-2 h-2 bg-purple-400 rounded-full animate-pulse" style="animation-delay: 0.5s"></div>
					<div class="w-2 h-2 bg-indigo-400 rounded-full animate-pulse" style="animation-delay: 1s"></div>
				</div>
			</div>
		</div>
	</footer>
</div>

<!-- Toast Notification -->
{#if showToast}
	<div class="fixed top-6 right-6 z-50 animate-in slide-in-from-right-full duration-300">
		<div class="bg-gradient-to-r from-green-500 to-emerald-500 text-white p-4 rounded-2xl shadow-2xl border border-white/20 backdrop-blur-sm">
			<div class="flex items-center space-x-3">
				<div class="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center">
					<span class="text-lg">✅</span>
				</div>
				<span class="font-medium">{toastMessage}</span>
			</div>
		</div>
	</div>
{/if}