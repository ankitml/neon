<script>
	let searchQuery = '';
	let searchResults = [];
	let isLoading = false;
	let error = '';
	let showToast = false;
	let toastMessage = '';

	async function handleSearch() {
		if (!searchQuery.trim()) return;
		
		isLoading = true;
		error = '';
		
		try {
			const response = await fetch(`http://localhost:8080/api/search?q=${encodeURIComponent(searchQuery)}`);
			
			if (!response.ok) {
				throw new Error('Search failed');
			}
			
			const data = await response.json();
			searchResults = data.results || [];
		} catch (err) {
			error = err.message || 'An error occurred';
			searchResults = [];
		} finally {
			isLoading = false;
		}
	}

	function copyToClipboard(text) {
		navigator.clipboard.writeText(text).then(() => {
			toastMessage = 'Quote copied to clipboard!';
			showToast = true;
			setTimeout(() => showToast = false, 3000);
		});
	}
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
	<main class="container mx-auto px-4 py-8 space-y-8">
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

		<!-- Search Results -->
		{#if searchResults.length > 0}
			<div class="space-y-6">
				<div class="flex items-center justify-between">
					<h2 class="text-2xl font-bold text-gray-800">Search Results</h2>
					<span class="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-medium">
						{searchResults.length} quotes found
					</span>
				</div>
				
				<div class="grid gap-4">
					{#each searchResults as quote}
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
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else if !isLoading && searchQuery}
			<div class="bg-white rounded-lg shadow-md p-8 text-center space-y-4">
				<div class="text-6xl">🔍</div>
				<h3 class="text-xl font-semibold text-gray-800">No quotes found</h3>
				<p class="text-gray-600">
					No quotes found for "<strong>{searchQuery}</strong>". 
					Try searching for themes like "courage", "love", or author names.
				</p>
			</div>
		{:else if !searchQuery}
			<div class="bg-white rounded-lg shadow-md p-8 text-center space-y-4">
				<div class="text-6xl">💭</div>
				<h3 class="text-xl font-semibold text-gray-800">Welcome to Quote Collections</h3>
				<p class="text-gray-600">
					Start by searching for quotes by theme, author, or keyword above.
					<br />
					Try searching for: <strong>life</strong>, <strong>success</strong>, <strong>Einstein</strong>, or <strong>courage</strong>
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