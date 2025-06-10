 package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Quote struct {
	ID       int      `json:"id"`
	Text     string   `json:"text"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type QuoteResponse struct {
	Quote    Quote  `json:"quote"`
	Total    int    `json:"total"`
	Category string `json:"category,omitempty"`
}

type TimeZoneInfo struct {
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Time         string `json:"time"`
	Timestamp    int64  `json:"timestamp"`
}

type TimeResponse struct {
	TimeZones []TimeZoneInfo `json:"timezones"`
	UpdatedAt string         `json:"updated_at"`
}

var quotes = []Quote{
	{1, "We don't make mistakes, just happy little accidents.", "Philosophy", []string{"mistakes", "learning", "positive"}},
	{2, "There's nothing wrong with having a tree as a friend.", "Nature", []string{"friendship", "nature", "trees"}},
	{3, "The secret to doing anything is believing that you can do it.", "Motivation", []string{"belief", "confidence", "success"}},
	{4, "Talent is a pursued interest. Anything that you're willing to practice, you can do.", "Learning", []string{"talent", "practice", "skill"}},
	{5, "You can do anything you want to do. This is your world.", "Empowerment", []string{"empowerment", "freedom", "possibilities"}},
	{6, "Go out on a limb — that's where the fruit is.", "Risk", []string{"courage", "risk", "opportunity"}},
	{7, "Look around. Look at what we have. Beauty is everywhere—you only have to look to see it.", "Beauty", []string{"beauty", "appreciation", "awareness"}},
	{8, "Just go out and talk to a tree. Make friends with it.", "Nature", []string{"nature", "friendship", "connection"}},
	{9, "There are no mistakes, only happy accidents.", "Philosophy", []string{"mistakes", "learning", "acceptance"}},
	{10, "All you need to paint is a few tools, a little instruction, and a vision in your mind.", "Art", []string{"art", "creativity", "vision"}},
	{11, "Clouds are very, very free.", "Nature", []string{"freedom", "clouds", "nature"}},
	{12, "In painting, you have unlimited power. You have the ability to move mountains.", "Art", []string{"power", "creativity", "imagination"}},
	{13, "You need the dark in order to show the light.", "Philosophy", []string{"contrast", "balance", "perspective"}},
	{14, "Just let go — and fall like a little waterfall.", "Philosophy", []string{"letting go", "flow", "peace"}},
	{15, "It's so important to do something every day that will make you happy.", "Happiness", []string{"happiness", "daily", "joy"}},
	{16, "Maybe this cloud has a little friend. Maybe this little friend's named Clyde.", "Friendship", []string{"friendship", "imagination", "clouds"}},
	{17, "I think there's an artist hidden at the bottom of every single one of us.", "Art", []string{"art", "potential", "creativity"}},
	{18, "We each see the world in our own way. That's what makes it such a special place.", "Philosophy", []string{"perspective", "uniqueness", "world"}},
	{19, "Anytime you learn, you gain.", "Learning", []string{"learning", "growth", "knowledge"}},
	{20, "Mix up a little more shadow color here, then we can put us a little shadow right in there.", "Art", []string{"technique", "shadow", "painting"}},
	{21, "Remember how free clouds are. They just lay around in the sky all day long.", "Nature", []string{"freedom", "clouds", "relaxation"}},
	{22, "If we're going to have animals around we all have to be concerned about them and take care of them.", "Nature", []string{"animals", "care", "responsibility"}},
	{23, "You can do anything here — the only prerequisite is that it makes you happy.", "Happiness", []string{"happiness", "freedom", "joy"}},
	{24, "Isn't it fantastic that you can change your mind and create all these happy things?", "Creativity", []string{"change", "creativity", "happiness"}},
	{25, "I guess I'm a little weird. I like to talk to trees and animals.", "Nature", []string{"nature", "quirky", "connection"}},
}

// US Time Zones
var usTimeZones = map[string]string{
	"Eastern":  "America/New_York",
	"Central":  "America/Chicago", 
	"Mountain": "America/Denver",
	"Pacific":  "America/Los_Angeles",
	"Alaska":   "America/Anchorage",
	"Hawaii":   "Pacific/Honolulu",
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Interactive Bob Ross Quotes with US Time Zones</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            font-family: 'Georgia', serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        
        .container {
            max-width: 1400px;
            margin: 0 auto;
            padding: 20px;
        }
        
        header {
            text-align: center;
            color: white;
            margin-bottom: 30px;
            padding: 20px;
        }
        
        h1 {
            font-size: 3em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        
        .subtitle {
            font-size: 1.2em;
            opacity: 0.9;
        }
        
        .main-grid {
            display: grid;
            grid-template-columns: 1fr 400px;
            gap: 30px;
            margin-bottom: 30px;
        }
        
        .left-panel {
            display: flex;
            flex-direction: column;
            gap: 30px;
        }
        
        .right-panel {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }
        
        .controls {
            background: white;
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
        }
        
        .control-group {
            display: flex;
            flex-wrap: wrap;
            gap: 15px;
            align-items: center;
            margin-bottom: 20px;
        }
        
        .search-box {
            flex: 1;
            min-width: 200px;
            padding: 12px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 25px;
            font-size: 16px;
            outline: none;
            transition: border-color 0.3s;
        }
        
        .search-box:focus {
            border-color: #667eea;
        }
        
        .category-select {
            padding: 12px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 25px;
            font-size: 16px;
            outline: none;
            background: white;
            cursor: pointer;
        }
        
        .btn {
            padding: 12px 25px;
            border: none;
            border-radius: 25px;
            font-size: 16px;
            cursor: pointer;
            transition: all 0.3s;
            font-weight: bold;
            text-decoration: none;
            display: inline-block;
        }
        
        .btn-primary {
            background: #667eea;
            color: white;
        }
        
        .btn-primary:hover {
            background: #5a6fd8;
            transform: translateY(-2px);
        }
        
        .btn-secondary {
            background: #764ba2;
            color: white;
        }
        
        .btn-secondary:hover {
            background: #6a4190;
            transform: translateY(-2px);
        }
        
        .quote-display {
            background: white;
            border-radius: 15px;
            padding: 40px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            text-align: center;
            position: relative;
            overflow: hidden;
        }
        
        .quote-display::before {
            content: '🎨';
            position: absolute;
            top: 20px;
            left: 30px;
            font-size: 2em;
            opacity: 0.1;
        }
        
        .quote-text {
            font-size: 1.8em;
            line-height: 1.6;
            margin-bottom: 20px;
            font-style: italic;
            color: #2c3e50;
        }
        
        .quote-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 15px;
            margin-top: 25px;
            padding-top: 20px;
            border-top: 2px solid #f0f0f0;
        }
        
        .quote-author {
            font-weight: bold;
            font-size: 1.1em;
            color: #667eea;
        }
        
        .quote-category {
            background: #667eea;
            color: white;
            padding: 5px 15px;
            border-radius: 20px;
            font-size: 0.9em;
        }
        
        .quote-tags {
            display: flex;
            gap: 8px;
            flex-wrap: wrap;
        }
        
        .tag {
            background: #f8f9fa;
            color: #6c757d;
            padding: 4px 12px;
            border-radius: 15px;
            font-size: 0.8em;
            border: 1px solid #e9ecef;
        }
        
        .timezones-panel {
            background: white;
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            height: fit-content;
        }
        
        .timezones-title {
            text-align: center;
            margin-bottom: 20px;
            color: #2c3e50;
            font-size: 1.4em;
            border-bottom: 2px solid #f0f0f0;
            padding-bottom: 15px;
        }
        
        .timezone-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 15px;
            margin-bottom: 10px;
            background: #f8f9fa;
            border-radius: 10px;
            border-left: 4px solid #667eea;
            transition: all 0.3s;
        }
        
        .timezone-item:hover {
            transform: translateX(5px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        
        .timezone-name {
            font-weight: bold;
            color: #2c3e50;
            font-size: 1.1em;
        }
        
        .timezone-abbr {
            font-size: 0.9em;
            color: #6c757d;
            margin-top: 2px;
        }
        
        .timezone-time {
            text-align: right;
        }
        
        .timezone-clock {
            font-family: 'Courier New', monospace;
            font-size: 1.2em;
            font-weight: bold;
            color: #667eea;
        }
        
        .timezone-date {
            font-size: 0.8em;
            color: #6c757d;
            margin-top: 2px;
        }
        
        .stats {
            background: white;
            border-radius: 15px;
            padding: 20px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            text-align: center;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 20px;
            margin-top: 15px;
        }
        
        .stat-item {
            padding: 15px;
            background: #f8f9fa;
            border-radius: 10px;
        }
        
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }
        
        .stat-label {
            color: #6c757d;
            font-size: 0.9em;
        }
        
        .loading {
            display: none;
            text-align: center;
            color: #667eea;
            font-size: 1.2em;
        }
        
        .error {
            background: #f8d7da;
            color: #721c24;
            padding: 15px;
            border-radius: 10px;
            margin: 20px 0;
            display: none;
        }
        
        .trees {
            position: fixed;
            bottom: 20px;
            right: 20px;
            font-size: 2em;
            opacity: 0.3;
            pointer-events: none;
        }
        
        .current-time-display {
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            padding: 15px;
            border-radius: 10px;
            margin-bottom: 20px;
            text-align: center;
        }
        
        .current-time {
            font-family: 'Courier New', monospace;
            font-size: 1.5em;
            font-weight: bold;
        }
        
        .live-indicator {
            display: inline-block;
            width: 8px;
            height: 8px;
            background: #2ecc71;
            border-radius: 50%;
            margin-right: 8px;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0% { opacity: 1; }
            50% { opacity: 0.5; }
            100% { opacity: 1; }
        }
        
        @media (max-width: 768px) {
            .container { padding: 10px; }
            h1 { font-size: 2em; }
            .quote-text { font-size: 1.4em; }
            .main-grid { 
                grid-template-columns: 1fr; 
                gap: 20px;
            }
            .control-group { flex-direction: column; }
            .search-box { min-width: 100%; }
            .timezone-item {
                flex-direction: column;
                text-align: center;
                gap: 10px;
            }
        }
        
        .fade-in {
            animation: fadeIn 0.5s ease-in;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🎨 Interactive Bob Ross Quotes</h1>
            <p class="subtitle">Learning Go through the wisdom of Bob Ross + Live US Time Zones</p>
        </header>
        
        <div class="main-grid">
            <div class="left-panel">
                <div class="controls">
                    <div class="control-group">
                        <input type="text" id="searchBox" class="search-box" placeholder="Search quotes by keyword...">
                        <select id="categorySelect" class="category-select">
                            <option value="">All Categories</option>
                            <option value="Philosophy">Philosophy</option>
                            <option value="Nature">Nature</option>
                            <option value="Art">Art</option>
                            <option value="Motivation">Motivation</option>
                            <option value="Learning">Learning</option>
                            <option value="Happiness">Happiness</option>
                            <option value="Friendship">Friendship</option>
                            <option value="Beauty">Beauty</option>
                            <option value="Empowerment">Empowerment</option>
                            <option value="Risk">Risk</option>
                            <option value="Creativity">Creativity</option>
                        </select>
                        <button class="btn btn-primary" onclick="getRandomQuote()">Random Quote</button>
                        <button class="btn btn-secondary" onclick="searchQuotes()">Search</button>
                    </div>
                </div>
                
                <div class="loading" id="loading">🌲 Finding your perfect quote...</div>
                <div class="error" id="error"></div>
                
                <div class="quote-display fade-in" id="quoteDisplay">
                    <div class="quote-text" id="quoteText">Click "Random Quote" to get started!</div>
                    <div class="quote-meta" id="quoteMeta" style="display: none;">
                        <div class="quote-author">— Bob Ross</div>
                        <div class="quote-category" id="quoteCategory"></div>
                        <div class="quote-tags" id="quoteTags"></div>
                    </div>
                </div>
                
                <div class="stats">
                    <h3>Quote Statistics</h3>
                    <div class="stats-grid">
                        <div class="stat-item">
                            <div class="stat-number" id="totalQuotes">{{.Total}}</div>
                            <div class="stat-label">Total Quotes</div>
                        </div>
                        <div class="stat-item">
                            <div class="stat-number" id="categories">{{.Categories}}</div>
                            <div class="stat-label">Categories</div>
                        </div>
                        <div class="stat-item">
                            <div class="stat-number" id="quotesViewed">0</div>
                            <div class="stat-label">Quotes Viewed</div>
                        </div>
                    </div>
                </div>
            </div>
            
            <div class="right-panel">
                <div class="timezones-panel">
                    <div class="current-time-display">
                        <div><span class="live-indicator"></span>Live US Time Zones</div>
                        <div class="current-time" id="currentTime"></div>
                    </div>
                    
                    <div class="timezones-title">🕐 US Time Zones</div>
                    <div id="timezonesContainer">
                        <!-- Time zones will be populated here -->
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <div class="trees">🌲🌲🌲</div>
    
    <script>
        let quotesViewed = 0;
        
        function showLoading() {
            document.getElementById('loading').style.display = 'block';
            document.getElementById('error').style.display = 'none';
        }
        
        function hideLoading() {
            document.getElementById('loading').style.display = 'none';
        }
        
        function showError(message) {
            document.getElementById('error').textContent = message;
            document.getElementById('error').style.display = 'block';
            hideLoading();
        }
        
        function updateQuoteDisplay(data) {
            const quote = data.quote;
            document.getElementById('quoteText').textContent = '"' + quote.text + '"';
            document.getElementById('quoteCategory').textContent = quote.category;
            
            const tagsContainer = document.getElementById('quoteTags');
            tagsContainer.innerHTML = '';
            quote.tags.forEach(tag => {
                const tagElement = document.createElement('span');
                tagElement.className = 'tag';
                tagElement.textContent = '#' + tag;
                tagsContainer.appendChild(tagElement);
            });
            
            document.getElementById('quoteMeta').style.display = 'flex';
            document.getElementById('quoteDisplay').classList.add('fade-in');
            
            quotesViewed++;
            document.getElementById('quotesViewed').textContent = quotesViewed;
            hideLoading();
        }
        
        function getRandomQuote() {
            showLoading();
            fetch('/api/random')
                .then(response => response.json())
                .then(data => updateQuoteDisplay(data))
                .catch(error => showError('Error fetching quote: ' + error.message));
        }
        
        function searchQuotes() {
            const searchTerm = document.getElementById('searchBox').value;
            const category = document.getElementById('categorySelect').value;
            
            if (!searchTerm && !category) {
                getRandomQuote();
                return;
            }
            
            showLoading();
            let url = '/api/search?';
            if (searchTerm) url += 'q=' + encodeURIComponent(searchTerm) + '&';
            if (category) url += 'category=' + encodeURIComponent(category);
            
            fetch(url)
                .then(response => response.json())
                .then(data => {
                    if (data.quote) {
                        updateQuoteDisplay(data);
                    } else {
                        showError('No quotes found matching your criteria.');
                    }
                })
                .catch(error => showError('Error searching quotes: ' + error.message));
        }
        
        function updateTimeZones() {
            fetch('/api/timezones')
                .then(response => response.json())
                .then(data => {
                    const container = document.getElementById('timezonesContainer');
                    container.innerHTML = '';
                    
                    data.timezones.forEach(tz => {
                        const tzElement = document.createElement('div');
                        tzElement.className = 'timezone-item';
                        
                        const date = new Date(tz.timestamp * 1000);
                        const timeStr = date.toLocaleTimeString('en-US', {
                            hour12: true,
                            hour: '2-digit',
                            minute: '2-digit',
                            second: '2-digit'
                        });
                        const dateStr = date.toLocaleDateString('en-US', {
                            month: 'short',
                            day: 'numeric',
                            year: 'numeric'
                        });
                        
                        tzElement.innerHTML = 
                            '<div>' +
                                '<div class="timezone-name">' + tz.name + '</div>' +
                                '<div class="timezone-abbr">' + tz.abbreviation + '</div>' +
                            '</div>' +
                            '<div class="timezone-time">' +
                                '<div class="timezone-clock">' + timeStr + '</div>' +
                                '<div class="timezone-date">' + dateStr + '</div>' +
                            '</div>';
                        
                        container.appendChild(tzElement);
                    });
                    
                    // Update current time display
                    const now = new Date();
                    document.getElementById('currentTime').textContent = now.toLocaleTimeString('en-US', {
                        hour12: true,
                        hour: '2-digit',
                        minute: '2-digit',
                        second: '2-digit'
                    });
                })
                .catch(error => console.error('Error fetching time zones:', error));
        }
        
        // Enter key support for search
        document.getElementById('searchBox').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                searchQuotes();
            }
        });
        
        // Initialize
        document.addEventListener('DOMContentLoaded', function() {
            setTimeout(getRandomQuote, 500);
            updateTimeZones();
            
            // Update time zones every second
            setInterval(updateTimeZones, 1000);
        });
    </script>
</body>
</html>
`

func getCategories() []string {
	categoryMap := make(map[string]bool)
	for _, quote := range quotes {
		categoryMap[quote.Category] = true
	}
	
	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}
	return categories
}

func getRandomQuote() Quote {
	rand.Seed(time.Now().UnixNano())
	return quotes[rand.Intn(len(quotes))]
}

func searchQuotes(searchTerm, category string) []Quote {
	var results []Quote
	searchLower := strings.ToLower(searchTerm)
	
	for _, quote := range quotes {
		matchesSearch := searchTerm == "" || 
			strings.Contains(strings.ToLower(quote.Text), searchLower) ||
			containsTag(quote.Tags, searchLower)
		
		matchesCategory := category == "" || quote.Category == category
		
		if matchesSearch && matchesCategory {
			results = append(results, quote)
		}
	}
	return results
}

func containsTag(tags []string, search string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), search) {
			return true
		}
	}
	return false
}

func getUSTimeZones() []TimeZoneInfo {
	var timeZones []TimeZoneInfo
	now := time.Now()
	
	for name, location := range usTimeZones {
		loc, err := time.LoadLocation(location)
		if err != nil {
			continue
		}
		
		localTime := now.In(loc)
		
		// Get timezone abbreviation
		abbr, _ := localTime.Zone()
		
		timeZones = append(timeZones, TimeZoneInfo{
			Name:         name,
			Abbreviation: abbr,
			Time:         localTime.Format("3:04:05 PM"),
			Timestamp:    localTime.Unix(),
		})
	}
	
	return timeZones
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	data := struct {
		Total      int
		Categories int
	}{
		Total:      len(quotes),
		Categories: len(getCategories()),
	}
	
	tmpl.Execute(w, data)
}

func randomQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	quote := getRandomQuote()
	response := QuoteResponse{
		Quote: quote,
		Total: len(quotes),
	}
	json.NewEncoder(w).Encode(response)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	searchTerm := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	
	results := searchQuotes(searchTerm, category)
	
	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No quotes found"})
		return
	}
	
	// Return a random quote from the results
	rand.Seed(time.Now().UnixNano())
	selectedQuote := results[rand.Intn(len(results))]
	
	response := QuoteResponse{
		Quote:    selectedQuote,
		Total:    len(results),
		Category: category,
	}
	json.NewEncoder(w).Encode(response)
}

func quoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}
	
	for _, quote := range quotes {
		if quote.ID == id {
			response := QuoteResponse{
				Quote: quote,
				Total: len(quotes),
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Quote not found"})
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories := getCategories()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categories,
		"total":      len(categories),
	})
}

func timeZonesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	timeZones := getUSTimeZones()
	response := TimeResponse{
		TimeZones: timeZones,
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Web routes
	http.HandleFunc("/", homeHandler)
	
	// API routes
	http.HandleFunc("/api/random", randomQuoteHandler)
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/quote", quoteByIDHandler)
	http.HandleFunc("/api/categories", categoriesHandler)
	http.HandleFunc("/api/timezones", timeZonesHandler)
	
	fmt.Println("🎨 Interactive Bob Ross Quotes + Time Zones Server Starting...")
	fmt.Println("🌲 Learning Go through the wisdom of Bob Ross")
	fmt.Println("🕐 Real-time US Time Zones with running seconds")
	fmt.Println("📡 Server running on http://localhost:8080")
	fmt.Println("🔍 Features:")
	fmt.Println("   • Interactive search and filtering")
	fmt.Println("   • Categories and tags")
	fmt.Println("   • Live US time zones (6 zones)")
	fmt.Println("   • RESTful API endpoints")
	fmt.Println("   • Real-time clock updates")
	fmt.Println("   • Responsive design")
	fmt.Println("🎯 Press Ctrl+C to stop the server")
	fmt.Println("")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}