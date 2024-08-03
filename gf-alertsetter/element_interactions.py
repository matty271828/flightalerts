# element_interactions.py
import time

from selenium.webdriver.common.by import By

def click_element(driver, xpath):
    """Clicks on the element specified by the given XPath."""
    wait_time = 1
    try:
        element = driver.find_element(By.XPATH, xpath)
        element.click()
        time.sleep(wait_time)  # Wait for the action to complete
        print(f"Element clicked successfully: {xpath}")
    except Exception as e:
        print(f"An error occurred while clicking element: {xpath}, Error: {e}")

def enter_text(driver, xpath, text):
    """Enters the given text into the element specified by the given XPath."""
    wait_time = 1
    try:
        element = driver.find_element(By.XPATH, xpath)
        element.clear()  # Clear any pre-existing text
        element.send_keys(text)
        time.sleep(wait_time)  # Wait for the action to complete
        print(f"Text '{text}' entered successfully in element: {xpath}")
    except Exception as e:
        print(f"An error occurred while entering text: {xpath}, Error: {e}")

def refresh_page(driver):
    wait_time = 5
    try: 
        print("Attempting to refresh the page")
        driver.refresh()
        time.sleep(wait_time)
    except Exception as e:
        print("An error occurred attempting to refresh the page")