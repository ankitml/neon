import json
import random

def generate_realistic_quotes(num_quotes):
    """
    Generates a list of synthetic but more realistic quotes in JSON format using sentence templates.
    """
    
    # Word lists for sentence construction
    nouns = ["river", "mountain", "tree", "ocean", "star", "sun", "moon", "world", "life", "time", "person", "hand", "eye", "work", "love", "hope", "fear", "joy", "sorrow"]
    adjectives = ["quick", "lazy", "bright", "dark", "happy", "sad", "beautiful", "ugly", "strong", "weak", "old", "young", "new", "great", "little"]
    verbs = ["jumps", "runs", "sleeps", "dreams", "thinks", "loves", "hates", "sees", "finds", "creates", "destroys", "inspires", "guides"]
    adverbs = ["quickly", "slowly", "happily", "sadly", "beautifully", "uglily", "strongly", "weakly", "truly", "falsely"]

    # Sentence templates
    templates = [
        "The {adjective} {noun} {adverb} {verb} the {adjective} {noun}.",
        "A {adjective} {noun} is like a {adjective} {noun}.",
        "To {verb} is to {adverb} understand the {noun}.",
        "Without {noun}, there is no {adjective} {noun}.",
        "The {noun} of {noun} is the {adjective} path to {noun}.",
        "{adverb}, the {noun} {verb}."
    ]

    tags = ["inspiration", "motivation", "life", "love", "wisdom", "humor", "philosophy", "truth", "friendship", "death", "hope", "failure", "success", "happiness", "knowledge", "science", "art", "imagination", "dreams"]
    categories = ["life", "love", "humor", "philosophy", "inspirational", "motivational", "general"]
    
    quotes_list = []
    for _ in range(num_quotes):
        template = random.choice(templates)
        # This makes sure that the same noun is not used twice in the same quote
        noun1 = random.choice(nouns)
        noun2 = random.choice(nouns)
        while noun1 == noun2:
            noun2 = random.choice(nouns)
            
        quote_text = template.format(
            noun=noun1,
            adjective=random.choice(adjectives),
            verb=random.choice(verbs),
            adverb=random.choice(adverbs)
        ).replace("{noun}", noun2, 1) # Replace the second {noun} with noun2
        
        quote_text = quote_text.capitalize()

        author_words = [random.choice(nouns).capitalize() for _ in range(2)]
        author_name = " ".join(author_words)
        
        num_tags = random.randint(2, 4)
        quote_tags = random.sample(tags, num_tags)
        
        popularity = round(random.random(), 4)
        
        category = random.choice(categories)
        
        quote_obj = {
            "Quote": quote_text,
            "Author": author_name,
            "Tags": quote_tags,
            "Popularity": popularity,
            "Category": category
        }
        quotes_list.append(quote_obj)
        
    return quotes_list

if __name__ == "__main__":
    number_of_quotes_to_generate = 1000
    
    quotes_data = generate_realistic_quotes(number_of_quotes_to_generate)
    
    output_filename = "realistic_quotes_generated.json"
    
    with open(output_filename, "w") as json_file:
        json.dump(quotes_data, json_file, indent=4)
        
    print(f"Successfully generated {number_of_quotes_to_generate} realistic quotes and saved them to '{output_filename}'")
