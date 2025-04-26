import json
import time
import undetected_chromedriver as uc
from bs4 import BeautifulSoup
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
import cffi
from bs4 import BeautifulSoup
ffi = cffi.FFI()
ffi.cdef("""
    char* get_products_links(char* url);
""")

def collect_product_info(driver, url):
    product_id = WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.XPATH, '//div[contains(text(), "Артикул: ")]'))).text.split('Артикул: ')[1]
    image_url = None
    page_source = str(driver.page_source)
    soup = BeautifulSoup(page_source, 'lxml')

    with open(f'product_{product_id}.html', 'w') as file:
        file.write(page_source)

    product_name = soup.find('div', attrs={"data-widget": 'webProductHeading'}).find(
        'h1').text.strip().replace('\t', '').replace('\n', ' ')
    try:
            ozon_card_price_element = soup.find(
                'span', string="c Ozon Картой").parent.find('div').find('span')
            product_ozon_card_price = ozon_card_price_element.text.strip(
            ) if ozon_card_price_element else ''

            price_element = soup.find(
                'span', string="без Ozon Карты").parent.parent.find('div').find_all('span')

            product_discount_price = price_element[0].text.strip(
            ) if price_element[0] else ''
            product_base_price = price_element[1].text.strip(
            ) if price_element[1] is not None else ''
    except:
        product_ozon_card_price = None
        product_discount_price = None


        # product price
        try:
            ozon_card_price_element = soup.find(
                'span', string="c Ozon Картой").parent.find('div').find('span')
        except AttributeError:
            card_price_div = soup.find(
                'div', attrs={"data-widget": "webPrice"}).findAll('span')


            product_discount_price = card_price_div[1].text.strip()


    try:
        img_element = WebDriverWait(driver, 10).until(EC.presence_of_element_located((By.CSS_SELECTOR, 'div[data-widget="webGallery"] img, img.zoom-image')))

        # Получаем URL маленького изображения
        small_image_url = img_element.get_attribute('src')
        # Преобразуем в URL большого изображения
        if small_image_url:
            # Заменяем параметры размера на максимальный
            image_url = (small_image_url.replace('/wc50/', '/wc1000/').replace('/wc100/', '/wc1000/').replace('/wc250/', '/wc1000/').split('?')[0])  # Удаляем параметры URL после ?
    except:
        image_url = None

    product_data = {
        'product_id': product_id,

        'product_image_url': image_url,
        'product_name': product_name,
        'product_ozon_card_price': product_ozon_card_price.replace("\u2009","").replace("\u2009","").replace("₽",""),
        'product_discount_price': product_discount_price.replace("\u2009","").replace("\u2009","").replace("₽",""),

    }
    return product_data
    get_products_links
@ffi.def_extern()
def get_products_links(url):
    driver = uc.Chrome()

    driver.get(ffi.string(url).decode('utf-8'))


     data = collect_product_info(driver, ffi.string(url).decode('utf-8'))
        json_data = json.dumps(data)

        driver.close()
        driver.quit()

        # Возвращаем копию строки, которую потом нужно освободить
        return ffi.new("char[]", json_data.encode('utf-8'))
if __name__ == '__main__':
    ffi.compile()