
"use client";

import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useState, useEffect } from "react";

export default function QuotesPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [results, setResults] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [facets, setFacets] = useState({});
  const [pagination, setPagination] = useState({});

  const [selectedCategories, setSelectedCategories] = useState([]);
  const [selectedTags, setSelectedTags] = useState([]);

  async function loadQuotes() {
    setIsLoading(true);
    setError("");

    try {
      const params = new URLSearchParams({
        page: "1",
        limit: "20",
        sort: "popularity",
        order: "desc",
        facets: "true",
      });

      if (searchQuery.trim()) {
        params.set("q", searchQuery.trim());
      }

      selectedCategories.forEach((cat) => params.append("categories[]", cat));
      selectedTags.forEach((tag) => params.append("tags[]", tag));

      const response = await fetch(
        `http://localhost:8080/api/search?${params}`
      );

      if (!response.ok) {
        throw new Error("Request failed");
      }

      const data = await response.json();
      setResults(data.quotes || data.results || []);
      setFacets(data.facets || {});
      setPagination(data.pagination || { total_count: data.count || 0 });
    } catch (err) {
      setError(err.message || "An error occurred");
      setResults([]);
    } finally {
      setIsLoading(false);
    }
  }

  function handleCategoryChange(category) {
    setSelectedCategories((prev) =>
      prev.includes(category)
        ? prev.filter((c) => c !== category)
        : [...prev, category]
    );
  }

  function handleTagChange(tag) {
    setSelectedTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : [...prev, tag]
    );
  }

  useEffect(() => {
    loadQuotes();
  }, [selectedCategories, selectedTags]);

  return (
    <div className="container mx-auto p-4">
      <div className="flex gap-8">
        <div className="w-1/4">
          <h2 className="text-xl font-bold mb-4">Filters</h2>
          <div>
            <h3 className="font-semibold mb-2">Categories</h3>
            {facets.categories &&
              facets.categories.map((category) => (
                <div key={category.value} className="flex items-center space-x-2">
                  <Checkbox
                    id={category.value}
                    checked={selectedCategories.includes(category.value)}
                    onCheckedChange={() => handleCategoryChange(category.value)}
                  />
                  <label htmlFor={category.value}>{category.value} ({category.count})</label>
                </div>
              ))}
          </div>
          <div className="mt-4">
            <h3 className="font-semibold mb-2">Tags</h3>
            {facets.tags &&
              facets.tags.map((tag) => (
                <div key={tag.value} className="flex items-center space-x-2">
                  <Checkbox
                    id={tag.value}
                    checked={selectedTags.includes(tag.value)}
                    onCheckedChange={() => handleTagChange(tag.value)}
                  />
                  <label htmlFor={tag.value}>{tag.value} ({tag.count})</label>
                </div>
              ))}
          </div>
        </div>
        <div className="w-3/4">
          <h1 className="text-3xl font-bold mb-4">Quotes</h1>
          <div className="flex w-full max-w-sm items-center space-x-2 mb-4">
            <Input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search for quotes..."
            />
            <Button onClick={loadQuotes}>
              Search
            </Button>
          </div>

          {isLoading && <p>Loading...</p>}
          {error && <p className="text-red-500">{error}</p>}

          <div className="grid gap-4">
            {results.map((quote, index) => (
              <Card key={index}>
                <CardContent className="pt-6">
                  <blockquote className="text-lg font-medium">
                    "{quote.quote}"
                  </blockquote>
                  <p className="text-right text-gray-500 mt-4">- {quote.author}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
