from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.service import Service
import time

def set_flight_alert(driver, origin, destination):
    # Find and click the origin input
    origin_input = driver.find_element(By.XPATH, '//input[@aria-label="Departure airport"]')
    origin_input.clear()
    origin_input.send_keys(origin)
    origin_input.send_keys(Keys.RETURN)

    time.sleep(2)  # Wait for suggestions to load

    # Find and click the destination input
    destination_input = driver.find_element(By.XPATH, '//input[@aria-label="Arrival airport"]')
    destination_input.clear()
    destination_input.send_keys(destination)
    destination_input.send_keys(Keys.RETURN)

    time.sleep(2)  # Wait for suggestions to load

    # Click the search button
    search_button = driver.find_element(By.XPATH, '//button[@aria-label="Search"]')
    search_button.click()

    time.sleep(5)  # Wait for search results to load

    # Click the track prices toggle
    track_button = driver.find_element(By.XPATH, '//button[@aria-label="Track prices"]')
    track_button.click()

    time.sleep(2)  # Wait for the toggle action to complete

def main():
    # Setup ChromeDriver path
    service = Service('path/to/chromedriver')
    driver = webdriver.Chrome(service=service)

    # Open Google Flights
    driver.get('https://www.google.com/flights')

    routes = [
        ('Manchester', 'Berlin'),
        ('Manchester', 'Paris'),
        ('Manchester', 'Lisbon'),
    ]

    for origin, destination in routes:
        set_flight_alert(driver, origin, destination)
        time.sleep(5)  # Wait a bit before setting up the next alert

    driver.quit()

if __name__ == "__main__":
    main()
