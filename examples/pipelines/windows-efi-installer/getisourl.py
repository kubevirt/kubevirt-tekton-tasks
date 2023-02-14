#!/usr/bin/python

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.select import Select

options = Options()
options.headless = True

driver = webdriver.Chrome(options=options)
driver.implicitly_wait(10)

try:
    driver.get("https://www.microsoft.com/en-us/software-download/windows11")

    select_object = Select(driver.find_element(By.ID, "product-edition"))
    select_object.select_by_index(1)

    button_element = driver.find_element(By.ID, "submit-product-edition")
    button_element.click()

    select_object = Select(driver.find_element(By.ID, "product-languages"))
    select_object.select_by_visible_text("English (United States)")

    button_element = driver.find_element(By.ID, "submit-sku")
    button_element.click()

    download_link = driver.find_element(By.LINK_TEXT, "64-bit Download")
    print(download_link.get_attribute("href"))
except:
    pass

driver.close()
