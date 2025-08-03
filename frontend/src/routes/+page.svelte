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
<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
	<!-- Header -->
	<header class="bg-blue-600 text-white shadow-lg">
		<div class="container mx-auto px-4 py-6">
			<h1 class="text-4xl font-bold text-center">📚 Themed Quote Collections</h1>
			<p class="text-center mt-2 opacity-90">Discover inspiring quotes by themes, authors, and keywords</p>
		</div>
	</header>

	<!-- Main Content -->
	<main class="container mx-auto px-4 py-8 space-y-6">
		<!-- Search Section -->
		<div class="bg-white rounded-lg shadow-md p-6">
			<form on:submit|preventDefault={handleSearch} class="space-y-4">
				<div class="space-y-2">
					<label class="block font-semibold text-gray-700" for="search">
						Search Quotes
					</label>
					<div class="flex gap-2">
						<input
							id="search"
							bind:value={searchQuery}
							type="text"
							placeholder="Enter themes like 'courage', 'innovation', or author names..."
							class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
						/>
						<button 
							type="submit" 
							disabled={isLoading}
							class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
						>
							{isLoading ? 'Searching...' : 'Search'}
						</button>
					</div>
				</div>
			</form>
		</div>

		<!-- Refinements Section -->
		<div class="bg-white rounded-lg shadow-md p-6">
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-semibold text-gray-800">Refinements</h3>
				<button 
					on:click={() => showRefinements = !showRefinements}
					class="px-4 py-2 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md transition-colors"
				>
					{showRefinements ? 'Hide' : 'Show'} Filters
				</button>
			</div>

			{#if showRefinements}
				<div class="space-y-6">
					<!-- Sorting Controls -->
					<div class="flex flex-wrap gap-4">
						<div class="flex items-center gap-2">
							<label for="sort-select" class="text-sm font-medium text-gray-700">Sort by:</label>
							<select 
								id="sort-select"
								bind:value={sort} 
								on:change={() => { page = 1; loadQuotes(); }}
								class="px-3 py-1 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
							>
								<option value="popularity">Popularity</option>
								<option value="created_at">Date Added</option>
								<option value="random">Random</option>
							</select>
						</div>
						<div class="flex items-center gap-2">
							<label for="order-select" class="text-sm font-medium text-gray-700">Order:</label>
							<select 
								id="order-select"
								bind:value={order} 
								on:change={() => { page = 1; loadQuotes(); }}
								class="px-3 py-1 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
							>
								<option value="desc">High to Low</option>
								<option value="asc">Low to High</option>
							</select>
						</div>
						{#if selectedCategories.length > 0 || selectedTags.length > 0}
							<button 
								on:click={clearAllFilters}
								class="px-3 py-1 text-sm bg-red-100 hover:bg-red-200 text-red-700 rounded-md transition-colors"
							>
								Clear All
							</button>
						{/if}
					</div>

					<!-- Active Filters -->
					{#if selectedCategories.length > 0 || selectedTags.length > 0}
						<div class="space-y-2">
							<h4 class="text-sm font-medium text-gray-700">Active Filters:</h4>
							<div class="flex flex-wrap gap-2">
								{#each selectedCategories as category}
									<span class="inline-flex items-center gap-1 bg-purple-100 text-purple-800 text-sm px-2 py-1 rounded-md">
										Category: {category}
										<button on:click={() => toggleCategory(category)} class="text-purple-600 hover:text-purple-800">✕</button>
									</span>
								{/each}
								{#each selectedTags as tag}
									<span class="inline-flex items-center gap-1 bg-blue-100 text-blue-800 text-sm px-2 py-1 rounded-md">
										Tag: {tag}
										<button on:click={() => toggleTag(tag)} class="text-blue-600 hover:text-blue-800">✕</button>
									</span>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Category Facets -->
					{#if facets.categories && facets.categories.length > 0}
						<div class="space-y-2">
							<h4 class="text-sm font-medium text-gray-700">Categories:</h4>
							<div class="flex flex-wrap gap-2">
								{#each facets.categories as category}
									<button 
										on:click={() => toggleCategory(category.value)}
										class="inline-flex items-center gap-1 px-3 py-1 text-sm rounded-md transition-colors {selectedCategories.includes(category.value) ? 'bg-purple-600 text-white' : 'bg-gray-100 hover:bg-gray-200 text-gray-700'}"
									>
										{category.value}
										<span class="text-xs opacity-75">({category.count})</span>
									</button>
								{/each}
							</div>
						</div>
					{/if}

					<!-- Tag Facets -->
					{#if facets.tags && facets.tags.length > 0}
						<div class="space-y-2">
							<h4 class="text-sm font-medium text-gray-700">Popular Tags:</h4>
							<div class="flex flex-wrap gap-2">
								{#each facets.tags as tag}
									<button 
										on:click={() => toggleTag(tag.value)}
										class="inline-flex items-center gap-1 px-3 py-1 text-sm rounded-md transition-colors {selectedTags.includes(tag.value) ? 'bg-blue-600 text-white' : 'bg-gray-100 hover:bg-gray-200 text-gray-700'}"
									>
										{tag.value}
										<span class="text-xs opacity-75">({tag.count})</span>
									</button>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<!-- Loading Bar -->
		{#if isLoading}
			<div class="w-full bg-gray-200 rounded-full h-2">
				<div class="bg-blue-600 h-2 rounded-full animate-pulse"></div>
			</div>
		{/if}

		<!-- Error Message -->
		{#if error}
			<div class="bg-red-50 border border-red-400 text-red-700 px-4 py-3 rounded-lg">
				❌ Error: {error}
			</div>
		{/if}

		<!-- Results -->
		{#if results.length > 0}
			<div class="space-y-6">
				<div class="flex items-center justify-between">
					<h2 class="text-2xl font-bold text-gray-800">
						{mode === 'search' ? 'Search Results' : 'Browse Quotes'}
					</h2>
					<span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
						{pagination.total_count || results.length} quotes found
						{#if mode === 'browse' && pagination.total_pages}
							• Page {pagination.page} of {pagination.total_pages}
						{/if}
					</span>
				</div>
				
				<div class="grid gap-4">
					{#each results as quote}
						<div class="bg-white rounded-lg shadow-md p-6 space-y-4 hover:shadow-lg transition-shadow">
							<blockquote class="text-lg font-medium text-gray-800 leading-relaxed">
								"{quote.quote}"
							</blockquote>
							
							<div class="flex items-center justify-between">
								<p class="text-gray-600 font-medium">— {quote.author}</p>
								<button 
									on:click={() => copyToClipboard(`"${quote.quote}" - ${quote.author}`)}
									class="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md transition-colors"
								>
									📋 Copy
								</button>
							</div>
							
							<div class="flex flex-wrap gap-2">
								{#if quote.category}
									<span class="inline-block bg-purple-100 text-purple-800 text-sm px-2 py-1 rounded-md">
										Category: {quote.category}
									</span>
								{/if}
								{#if quote.tags && quote.tags.length > 0}
									{#each quote.tags as tag}
										<span class="inline-block bg-gray-100 text-gray-700 text-sm px-2 py-1 rounded-md">
											{tag}
										</span>
									{/each}
								{/if}
								{#if quote.relevance}
									<span class="inline-block bg-green-100 text-green-800 text-sm px-2 py-1 rounded-md">
										Score: {quote.relevance.toFixed(2)}
									</span>
								{/if}
								{#if quote.popularity}
									<span class="inline-block bg-blue-100 text-blue-800 text-sm px-2 py-1 rounded-md">
										Popularity: {quote.popularity.toFixed(3)}
									</span>
								{/if}
							</div>
						</div>
					{/each}
				</div>

				<!-- Pagination -->
				{#if mode === 'browse' && pagination.total_pages > 1}
					<div class="flex justify-center items-center gap-2 mt-8">
						<button 
							on:click={() => changePage(pagination.page - 1)}
							disabled={!pagination.has_prev}
							class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
						>
							← Previous
						</button>
						
						<span class="px-4 py-2 text-sm text-gray-600">
							Page {pagination.page} of {pagination.total_pages}
						</span>
						
						<button 
							on:click={() => changePage(pagination.page + 1)}
							disabled={!pagination.has_next}
							class="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
						>
							Next →
						</button>
					</div>
				{/if}
			</div>
		{:else if !isLoading && searchQuery && mode === 'search'}
			<div class="bg-white rounded-lg shadow-md p-8 text-center space-y-4">
				<div class="text-6xl">🔍</div>
				<h3 class="text-xl font-semibold text-gray-800">No quotes found</h3>
				<p class="text-gray-600">
					No quotes found for "<strong>{searchQuery}</strong>". 
					Try searching for themes like "courage", "love", or author names.
				</p>
			</div>
		{:else if !isLoading && mode === 'browse' && results.length === 0}
			<div class="bg-white rounded-lg shadow-md p-8 text-center space-y-4">
				<div class="text-6xl">📂</div>
				<h3 class="text-xl font-semibold text-gray-800">No quotes match your filters</h3>
				<p class="text-gray-600">
					Try adjusting your category or tag filters to see more results.
				</p>
			</div>
		{/if}
	</main>

	<!-- Footer -->
	<footer class="mt-16 py-8 bg-gray-50 text-center">
		<p class="text-gray-500">
			Powered by ParadeDB full-text search • Built with SvelteKit & Tailwind CSS
		</p>
	</footer>
</div>

<!-- Toast Notification -->
{#if showToast}
	<div class="fixed top-4 right-4 bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg z-50">
		{toastMessage}
	</div>
{/if}