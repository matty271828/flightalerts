# Flight Alerts Automation

This project automates the process of setting up flight alerts on Google Flights using Selenium.

## Prerequisites

- Python 3.x
- Google Chrome
- ChromeDriver

## Installation

1. **Clone the repository**:

    ```bash
    git clone https://github.com/matty271828/flight_alerts
    cd flight_alerts
    ```

2. **Install dependencies**:

    ```bash
    pip install -r requirements.txt
    ```

3. **Download ChromeDriver**:

    Download ChromeDriver from [ChromeDriver](https://sites.google.com/chromium.org/driver/) and place it in a directory included in your system's PATH or specify its location in `main.py`.

## Usage

1. **Edit `main.py`**:

    Update the `routes` list with the origin and destination pairs for which you want to set up flight alerts.

    ```python
    routes = [
        ('New York', 'Los Angeles'),
        ('Chicago', 'Miami'),
        ('San Francisco', 'Seattle'),
        # Add more routes as needed
    ]
    ```

2. **Run the script**:

    ```bash
    python main.py
    ```

The script will open Google Flights, search for each route, and set up price alerts.

