from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.service import Service
import time

def set_flight_alert(driver, origin, destination):
    time.sleep(2)  # Wait for the page to load
    try:
        # Click the cookies "Accept All" button
        accept_cookies_xpath = '//*[@id="yDmH0d"]/c-wiz/div/div/div/div[2]/div[1]/div[3]/div[1]/div[1]/form[2]/div/div/button/span'
        accept_cookies_button = driver.find_element(By.XPATH, accept_cookies_xpath)
        accept_cookies_button.click()
        time.sleep(1)  # Wait for the action to complete
        print("Cookies 'Accept All' button clicked successfully.")
        
        # Enter the destination
        destination_input_xpath = '//*[@id="i23"]/div[4]/div/div/div[1]/div/div/input'
        destination_input = driver.find_element(By.XPATH, destination_input_xpath)
        destination_input.clear()  # Clear any pre-existing text
        destination_input.send_keys(destination)
        time.sleep(1)  # Wait for the action to complete
        print(f"Destination '{destination}' entered successfully.")
        
        time.sleep(1)  # Wait for input to load
        
        search_button_xpath = '//*[@id="c2"]'  # XPath for the search button
        search_button = driver.find_element(By.XPATH, search_button_xpath)
        search_button.click()
        time.sleep(1)  # Wait for the action to complete
        print("Search button clicked successfully.")

        # Enter the origin
        origin_input_xpath = '//*[@id="i23"]/div[1]/div/div/div[1]/div/div/input'
        origin_input = driver.find_element(By.XPATH, origin_input_xpath)
        origin_input.clear()  # Clear any pre-existing text
        origin_input.send_keys(origin)
        time.sleep(1)  # Wait for the action to complete
        print(f"Origin '{origin}' entered successfully.")
        
        time.sleep(2)  # Wait for input to load
        search_button_xpath = '//*[@id="c111"]'  # XPath for the search button
        search_button = driver.find_element(By.XPATH, search_button_xpath)
        search_button.click()
        time.sleep(1)  # Wait for the action to complete
        print("Search button clicked successfully.")
    
    except Exception as e:
        print(f"An error occurred: {e}")

def main():
    # Setup ChromeDriver path
    service = Service('/opt/homebrew/Caskroom/chromedriver/127.0.6533.72/chromedriver-mac-arm64/chromedriver')
    driver = webdriver.Chrome(service=service)

    # Open Google Flights
    driver.get('https://www.google.com/travel/flights?gl=GB&hl=en-GB')

    routes = [('Manchester', 'Berlin')]

    for origin, destination in routes:
        set_flight_alert(driver, origin, destination)
        time.sleep(500)  # Wait a bit before setting up the next alert

    driver.quit()

if __name__ == "__main__":
    main()
