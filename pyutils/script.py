import cv2
import numpy as np
import os

def is_valid(contour):
    x, y, w, h = cv2.boundingRect(contour)
    return 200 <= w <= 250 and 200 <= h <= 250

def run(fname, start=0):
    # Load the image
    image = cv2.imread(fname)
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

    # Threshold the image to create a binary image
    _, binary = cv2.threshold(gray, 240, 255, cv2.THRESH_BINARY_INV)

    # Find contours
    contours, _ = cv2.findContours(binary, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    contours = [
        c
        for c in contours
        if is_valid(c)
    ]
    grid_width = 4
    grid_height = len(contours) / grid_width
    min_x = min(cv2.boundingRect(c)[0] for c in contours)
    max_x = max(cv2.boundingRect(c)[0] for c in contours)
    min_y = min(cv2.boundingRect(c)[1] for c in contours)
    max_y = max(cv2.boundingRect(c)[1] for c in contours)

    # Sort contours by position (optional)
    def sort_key(countour):
        raw_x, raw_y, w, h = cv2.boundingRect(countour)

        x = round((raw_x - min_x) * grid_width / (max_x - min_x))
        y = round((raw_y - min_y) * grid_height / (max_y - min_y))

        return (y, x)

    contours = sorted(contours, key=sort_key)

    # Extract each thumbnail
    for n, cnt in enumerate(contours, start=start + 1):
        x, y, w, h = cv2.boundingRect(cnt)
        thumbnail = image[y:y+h, x:x+w]
        cv2.imwrite(os.path.join('.', f'piece{n}.png'), thumbnail)


run(fname='/Users/rsalmeidafl/Downloads/geometrix1.png', start=0)
run(fname='/Users/rsalmeidafl/Downloads/geometrix2.png', start=16)
