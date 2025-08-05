import json
import random

class QuoteGenerator:
    """
    A sophisticated quote generator that uses thematic word banks and complex sentence
    structures to create realistic and thematically coherent quotes.
    """

    def __init__(self):
        self.themes = {
            "nature": {
                "nouns": ["river", "mountain", "tree", "ocean", "star", "sun", "moon", "flower", "forest", "wind"],
                "adjectives": ["serene", "vast", "ancient", "whispering", "golden", "silent", "wild", "gentle"],
                "verbs": ["flows", "climbs", "grows", "dreams", "shines", "wanders", "breathes", "dances"],
                "adverbs": ["softly", "endlessly", "quietly", "majestically", "gently", "wildly"]
            },
            "wisdom": {
                "nouns": ["mind", "truth", "knowledge", "silence", "ignorance", "path", "journey", "question", "answer"],
                "adjectives": ["profound", "simple", "hidden", "eternal", "fleeting", "unseen", "clear"],
                "verbs": ["reveals", "hides", "understands", "questions", "guides", "illuminates", "teaches"],
                "adverbs": ["truly", "deeply", "patiently", "silently", "honestly", "wisely"]
            },
            "love": {
                "nouns": ["heart", "soul", "moment", "touch", "glance", "promise", "memory", "fire"],
                "adjectives": ["tender", "fierce", "unspoken", "unbreakable", "gentle", "passionate", "true"],
                "verbs": ["burns", "longs", "heals", "connects", "remembers", "forgives", "endures"],
                "adverbs": ["forever", "deeply", "softly", "unconditionally", "truly", "passionately"]
            },
            "ambition": {
                "nouns": ["dream", "goal", "summit", "horizon", "fire", "will", "struggle", "victory"],
                "adjectives": ["unrelenting", "bold", "daring", "solitary", "grand", "fierce", "audacious"],
                "verbs": ["climbs", "strives", "achieves", "conquers", "dares", "builds", "forges"],
                "adverbs": ["boldly", "relentlessly", "fearlessly", "tirelessly", "steadfastly", "audaciously"]
            },
            # --- New Eastern Spirituality Themes ---
            "zen": {
                "nouns": ["mind", "moment", "nothingness", "koan", "garden", "breath", "emptiness", "ego"],
                "adjectives": ["empty", "still", "direct", "unbound", "simple", "formless", "clear"],
                "verbs": ["sits", "observes", "accepts", "releases", "calms", "awakens", "points"],
                "adverbs": ["simply", "directly", "intently", "calmly", "presently", "effortlessly"]
            },
            "taoism": {
                "nouns": ["tao", "flow", "river", "valley", "balance", "yin", "yang", "way"],
                "adjectives": ["effortless", "yielding", "spontaneous", "harmonious", "natural", "unassuming"],
                "verbs": ["flows", "yields", "balances", "adapts", "harmonizes", "embraces", "lets-go"],
                "adverbs": ["effortlessly", "spontaneously", "naturally", "gently", "harmoniously"]
            },
            "mindfulness": {
                "nouns": ["breath", "sensation", "thought", "moment", "awareness", "presence", "body", "feeling"],
                "adjectives": ["present", "aware", "fleeting", "observant", "non-judgmental", "calm", "grounded"],
                "verbs": ["notices", "breathes", "grounds", "observes", "anchors", "accepts", "feels"],
                "adverbs": ["mindfully", "presently", "calmly", "attentively", "gently"]
            },
            "karma": {
                "nouns": ["action", "intention", "consequence", "seed", "fruit", "cycle", "cause", "effect"],
                "adjectives": ["skillful", "unskillful", "intentional", "inevitable", "subtle", "ripening"],
                "verbs": ["plants", "ripens", "returns", "creates", "shapes", "determines", "follows"],
                "adverbs": ["inevitably", "skillfully", "intentionally", "unfailingly"]
            },
            "impermanence": {
                "nouns": ["moment", "cloud", "river", "wave", "form", "feeling", "life", "breath"],
                "adjectives": ["fleeting", "transient", "changing", "ephemeral", "impermanent", "shifting"],
                "verbs": ["arises", "passes", "changes", "fades", "flows", "shifts", "dissolves"],
                "adverbs": ["constantly", "inevitably", "fleetingly", "momentarily"]
            }
        }

        self.templates = [
            "The {adjective} {noun} {adverb} {verb} the {adjective} {noun}.",
            "In the {noun} of a {adjective} {noun}, one {adverb} {verb}.",
            "To {verb} is to understand that every {noun} is a {adjective} {noun}.",
            "The {adjective} {noun} does not {verb}; it {adverb} {verb}.",
            "Without {noun}, a {noun} is but a {adjective} {noun}.",
            "Seek the {noun}, and you will {verb} the {adjective} {noun}.",
            "{adverb}, the {noun} {verb} beyond the {adjective} {noun}."
        ]

        self.first_names = ["Aethel", "Elara", "Kael", "Thorne", "Seraphina", "Orion", "Lyra", "Caius", "Rhiannon", "Jaxon", "Ananda", "Bodhi", "Kensho", "Satori"]
        self.last_names = ["Blackwood", "Stonehaven", "Ironhand", "Silverbow", "Nightwind", "Shadowend", "Starfall", "Fireheart", "Dharma", "Sangha"]

    def _generate_author(self):
        """Creates a more realistic-sounding author name."""
        return f"{random.choice(self.first_names)} {random.choice(self.last_names)}"

    def generate_quotes(self, num_quotes):
        """Generates a list of synthetic but more realistic quotes."""
        quotes_list = []
        for _ in range(num_quotes):
            theme_name = random.choice(list(self.themes.keys()))
            theme = self.themes[theme_name]

            template = random.choice(self.templates)

            # Ensure unique words for placeholders
            words = {
                "noun": random.sample(theme["nouns"], 2) if len(theme["nouns"]) > 1 else theme["nouns"] * 2,
                "adjective": random.sample(theme["adjectives"], 2) if len(theme["adjectives"]) > 1 else theme["adjectives"] * 2,
                "verb": random.sample(theme["verbs"], 2) if len(theme["verbs"]) > 1 else theme["verbs"] * 2,
                "adverb": random.sample(theme["adverbs"], 2) if len(theme["adverbs"]) > 1 else theme["adverbs"] * 2,
            }

            # Fill the template with unique words from the chosen theme
            quote_text = template.format(
                noun=words["noun"][0],
                adjective=words["adjective"][0],
                verb=words["verb"][0],
                adverb=words["adverb"][0]
            ).replace("{noun}", words["noun"][1], 1).replace("{adjective}", words["adjective"][1], 1).replace("{verb}", words["verb"][1], 1)

            quote_text = quote_text.capitalize()

            author = self._generate_author()

            num_tags = random.randint(2, 4)
            tags = random.sample(list(self.themes.keys()), k=min(num_tags, len(self.themes))) # Sample from theme keys for tags
            if theme_name not in tags:
                tags.append(theme_name)

            # Popularity can be influenced by theme
            base_popularity = 0.6 if theme_name in ["zen", "taoism", "mindfulness", "love", "wisdom"] else 0.4
            popularity = round(base_popularity + random.uniform(-0.2, 0.2), 4)

            quote_obj = {
                "Quote": quote_text,
                "Author": author,
                "Tags": tags,
                "Popularity": max(0, min(1, popularity)),  # Ensure popularity is between 0 and 1
                "Category": theme_name
            }
            quotes_list.append(quote_obj)

        return quotes_list

if __name__ == "__main__":
    generator = QuoteGenerator()

    number_of_quotes_to_generate = 1000
    quotes_data = generator.generate_quotes(number_of_quotes_to_generate)

    output_filename = "sophisticated_quotes_generated.json"

    with open(output_filename, "w") as json_file:
        json.dump(quotes_data, json_file, indent=4)

    print(f"Successfully generated {number_of_quotes_to_generate} sophisticated quotes and saved them to '{output_filename}'")
