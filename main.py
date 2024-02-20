import qrcode
import os
import math
from base45 import b45encode, b45decode


def getSize(filename):
    return os.path.getsize(filename)



import cv2
import os
from PIL import Image
from pyzbar.pyzbar import decode


def video_to_images(video_path, output_dir):
    os.makedirs(output_dir, exist_ok=True)
    cap = cv2.VideoCapture(video_path)
    if not cap.isOpened():
        print("Error opening video:", video_path)
        return

    frame_rate = cap.get(cv2.CAP_PROP_FPS)

    frame_count = 0
    while True:
        ret, frame = cap.read()
        if not ret:
            print("No more frames available")
            break

        output_path = "{}/frame_{:09}.jpg".format(output_dir, frame_count)

        cv2.imwrite(output_path, frame)

        if frame_count % 100 == 0:
            print(f"Extracted {frame_count} frames")

        frame_count += 1

    cap.release()

    print(f"Extracted a total of {frame_count} frames.")


def split_image(image_path, output_dir, target_size=(544, 544)):
    # Load the image
    img = Image.open(image_path)

    # Get image dimensions
    width, height = img.size

    # Calculate the number of rows and columns to fit the target size
    rows = max(1, height // target_size[1])
    cols = max(1, width // target_size[0])

    # Define coordinates for each sub-image
    sub_images = []
    for i in range(rows):
        for j in range(cols):
            # Calculate sub-image coordinates based on target size and position
            x1 = j * target_size[0]
            y1 = i * target_size[1]
            x2 = min(x1 + target_size[0], width)
            y2 = min(y1 + target_size[1], height)

            # Crop the sub-image
            sub_image = img.crop((x1, y1, x2, y2))

            # Resize the sub-image if necessary
            if sub_image.size != target_size:
                sub_image = sub_image.resize(target_size)

            sub_images.append(sub_image)

    # Save each sub-image with a sequential number
    for i, sub_image in enumerate(sub_images):
        output_path = f"{output_dir}/{os.path.splitext(os.path.basename(image_path))[0]}_{i + 1}.jpg"
        sub_image.save(output_path)


# # Example usage
video_path = "111.mp4"
output_dir = "temp2"

video_to_images(video_path, output_dir)



image_dir = "./temp2"
output_dir = "./t2"
os.makedirs(output_dir, exist_ok=True)

for filename in os.listdir(image_dir):
    if filename.lower().endswith((".jpg", ".jpeg", ".png")):
        image_path = os.path.join(image_dir, filename)
        split_image(image_path, output_dir)

print("Images split successfully!")


totalErrs = 0

dir = "./t2/"

file = b''
for filename in sorted(os.listdir(dir)):
    if os.path.isfile(os.path.join(dir, filename)):
        print(filename, "===>")
        data = decode(Image.open(dir + filename))
        try:
            print(data[0].data)
            file += b45decode(data[0].data)
        except:
            totalErrs += 1
            print("QR IS CORRUPTED")
            print(data)

print("TOTAL ERRORS ", totalErrs)
with open("my_file55555.pdf", "wb") as binary_file:
    # Write bytes to file
    binary_file.write(file)
