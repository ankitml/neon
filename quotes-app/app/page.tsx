"use client"

import { useState, useEffect } from "react"
import { Search, Filter, Heart, Copy, Share2, BookOpen, User, Calendar, TrendingUp } from "lucide-react"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useToast } from "@/hooks/use-toast"

interface Quote {
  id: string
  quote: string
  author: string
  category: string
  tags: string[]
}

const sortOptions = [
  { value: "popularity", label: "Most Popular" },
  { value: "recent", label: "Most Recent" },
  { value: "author", label: "Author A-Z" },
  { value: "length", label: "Shortest First" },
]

export default function QuotesApp() {
  const [searchQuery, setSearchQuery] = useState("")
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [sortBy, setSortBy] = useState("popularity")
  const [quotes, setQuotes] = useState<Quote[]>([])
  const [facets, setFacets] = useState<any>({})
  const [isLoading, setIsLoading] = useState(false)
  const { toast } = useToast()

  const searchQuotes = async (query: string, categories: string[], tags: string[], sort: string) => {
    setIsLoading(true)
    try {
      const params = new URLSearchParams({
        page: "1",
        limit: "20",
        sort: sort,
        order: "desc",
        facets: "true",
      });

      if (query.trim()) {
        params.set("q", query.trim());
      }

      categories.forEach((cat) => params.append("categories[]", cat));
      tags.forEach((tag) => params.append("tags[]", tag));

      const response = await fetch(
        `http://localhost:8080/api/search?${params}`
      );

      if (!response.ok) {
        throw new Error("Request failed");
      }

      const data = await response.json();
      setQuotes(data.quotes || data.results || []);
      setFacets(data.facets || {});
    } catch (err) {
      toast({
        title: "Error",
        description: "Failed to fetch quotes.",
        variant: "destructive",
      })
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    searchQuotes(searchQuery, selectedCategories, selectedTags, sortBy)
  }, [searchQuery, selectedCategories, selectedTags, sortBy])

  const handleCategoryChange = (category: string) => {
    setSelectedCategories((prev) =>
      prev.includes(category)
        ? prev.filter((c) => c !== category)
        : [...prev, category]
    );
  };

  const handleTagChange = (tag: string) => {
    setSelectedTags((prev) =>
      prev.includes(tag)
        ? prev.filter((t) => t !== tag)
        : [...prev, tag]
    );
  };

  const handleCopyQuote = (quote: Quote) => {
    navigator.clipboard.writeText(`"${quote.quote}" - ${quote.author}`)
    toast({
      title: "Quote copied!",
      description: "The quote has been copied to your clipboard.",
    })
  }

  const handleShareQuote = (quote: Quote) => {
    if (navigator.share) {
      navigator.share({
        title: `Quote by ${quote.author}`,
        text: `"${quote.quote}" - ${quote.author}`,
      })
    } else {
      handleCopyQuote(quote)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
      {/* Header with Search */}
      <div className="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-slate-200">
        <div className="max-w-6xl mx-auto px-4 py-4">
          <div className="flex items-center gap-4 mb-4">
            <div className="flex items-center gap-2">
              <BookOpen className="w-8 h-8 text-slate-700" />
              <h1 className="text-2xl font-bold text-slate-800">QuoteVault</h1>
            </div>
          </div>

          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-slate-400 w-5 h-5" />
            <Input
              placeholder="Search quotes, authors, or topics..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10 h-12 text-lg border-slate-300 focus:border-slate-500 focus:ring-slate-500"
            />
          </div>
        </div>
      </div>

      <div className="max-w-6xl mx-auto px-4 py-6">
        {/* Refinements Section */}
        <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 mb-8">
          <div className="flex items-center gap-2 mb-4">
            <Filter className="w-5 h-5 text-slate-600" />
            <h2 className="text-lg font-semibold text-slate-800">Refine Results</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Categories */}
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">Category</label>
              <div className="flex flex-wrap gap-2">
                {facets.categories?.map((category: any) => (
                  <Button
                    key={category.value}
                    variant={selectedCategories.includes(category.value) ? "default" : "outline"}
                    size="sm"
                    onClick={() => handleCategoryChange(category.value)}
                    className="text-xs"
                  >
                    {category.value} ({category.count})
                  </Button>
                ))}
              </div>
            </div>
            
            {/* Tags */}
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">Tags</label>
              <div className="flex flex-wrap gap-2">
                {facets.tags?.map((tag: any) => (
                  <Button
                    key={tag.value}
                    variant={selectedTags.includes(tag.value) ? "default" : "outline"}
                    size="sm"
                    onClick={() => handleTagChange(tag.value)}
                    className="text-xs"
                  >
                    {tag.value} ({tag.count})
                  </Button>
                ))}
              </div>
            </div>

            {/* Sort Options */}
            <div>
              <label className="block text-sm font-medium text-slate-700 mb-2">Sort By</label>
              <Select value={sortBy} onValueChange={setSortBy}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {sortOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {/* Results Count */}
            <div className="flex items-end">
              <div className="text-sm text-slate-600">
                <span className="font-semibold">{quotes.length}</span> quotes found
              </div>
            </div>
          </div>
        </div>

        {/* Results Section */}
        <div className="space-y-6">
          {searchQuery === "" && (
            <div className="flex items-center gap-2 mb-6">
              <TrendingUp className="w-5 h-5 text-slate-600" />
              <h2 className="text-xl font-semibold text-slate-800">Popular Quotes</h2>
            </div>
          )}

          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {[...Array(6)].map((_, i) => (
                <Card key={i} className="animate-pulse">
                  <CardContent className="p-6">
                    <div className="h-4 bg-slate-200 rounded mb-4"></div>
                    <div className="h-4 bg-slate-200 rounded mb-4 w-3/4"></div>
                    <div className="h-3 bg-slate-200 rounded w-1/2"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {quotes.map((quote) => (
                <Card
                  key={quote.id}
                  className="group hover:shadow-lg transition-all duration-300 border-slate-200 hover:border-slate-300"
                >
                  <CardContent className="p-6">
                    <blockquote className="text-lg text-slate-700 mb-4 leading-relaxed">"{quote.quote}"</blockquote>

                    <div className="flex items-center gap-2 mb-4">
                      <User className="w-4 h-4 text-slate-500" />
                      <cite className="text-slate-600 font-medium not-italic">{quote.author}</cite>
                    </div>

                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Badge variant="secondary" className="text-xs">
                          {quote.category}
                        </Badge>
                      </div>

                      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                        <Button variant="ghost" size="icon" className="w-8 h-8" onClick={() => handleCopyQuote(quote)}>
                          <Copy className="w-4 h-4" />
                        </Button>
                        <Button variant="ghost" size="icon" className="w-8 h-8" onClick={() => handleShareQuote(quote)}>
                          <Share2 className="w-4 h-4" />
                        </Button>
                        <Button variant="ghost" size="icon" className="w-8 h-8">
                          <Heart className="w-4 h-4" />
                        </Button>
                      </div>
                    </div>

                    <div className="flex flex-wrap gap-1 mt-3">
                      {quote.tags?.map((tag:any) => (
                        <Badge key={tag} variant="outline" className="text-xs">
                          #{tag}
                        </Badge>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}

          {quotes.length === 0 && !isLoading && (
            <div className="text-center py-12">
              <BookOpen className="w-16 h-16 text-slate-300 mx-auto mb-4" />
              <h3 className="text-xl font-semibold text-slate-600 mb-2">No quotes found</h3>
              <p className="text-slate-500">Try adjusting your search terms or filters</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
