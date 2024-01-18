#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include <cuda_runtime.h>
#include <cuda.h>
#include "helper_cuda.h"

#define STB_IMAGE_IMPLEMENTATION
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image.h"
#include "stb_image_write.h"

#define COLOR_CHANNELS 1

__global__ void marginirajGPE(const unsigned char *imageIn, unsigned char *imageOut, const int width, const int height)
{
    
// TODO #1.2.1 Določi pixel, ki ga bo nit obdelala
    int i = blockIdx.y * blockDim.y + threadIdx.y;
    int j = blockIdx.x * blockDim.x + threadIdx.x;

    int levi_rob = 0;
    int desni_rob = width - 1;
    int zgornji_rob = 0;
    int spodnji_rob = height - 1;

// TODO #1.2.2 Določi vrednost vsem sosedom
    int zgoraj = 0;
    int spodaj = 0;
    int levo = 0;
    int desno = 0;
    int zgoraj_levo = 0;
    int zgoraj_desno = 0;
    int spodaj_levo = 0;
    int spodaj_desno = 0;

    if (i != zgornji_rob) {
        zgoraj = imageIn[(i - 1) * width + j];
    }
    if (i != spodnji_rob) {
        spodaj = imageIn[(i + 1) * width + j];
    }
    if (j != levi_rob) {
        levo = imageIn[i * width + j - 1];
    }
    if (j != desni_rob) {
        desno = imageIn[i * width + j + 1];
    }
    if (i != zgornji_rob && j != levi_rob) {
        zgoraj_levo = imageIn[(i - 1) * width + j - 1];
    }
    if (i != zgornji_rob && j != desni_rob) {
        zgoraj_desno = imageIn[(i - 1) * width + j + 1];
    }
    if (i != spodnji_rob && j != levi_rob) {
        spodaj_levo = imageIn[(i + 1) * width + j - 1];
    }
    if (i != spodnji_rob && j != desni_rob) {
        spodaj_desno = imageIn[(i + 1) * width + j + 1];
    }

// TODO #1.2.3 Določi vrednost piksla glede na vrednosti sosedov
    int gx = - zgoraj_levo - 2 * levo - spodaj_levo + zgoraj_desno + 2 * desno + spodaj_desno;
    int gy = zgoraj_levo + 2 * zgoraj + zgoraj_desno - spodaj_levo - 2 * spodaj - spodaj_desno;

    float g = sqrt((float)(gx * gx + gy * gy));

    if (g > 255) {
        g = 255;
    }

    imageOut[i * width + j] = g;
}

void marginirajCPE(const unsigned char *imageIn, unsigned char *imageOut, const int width, const int height)
{
    printf("CPE DELA K ZMEŠAN\n");

// TODO #2.2.1 Implementiraj metodo za izvajanje na gostitelju
// TODO #2.2.2 Sprehod po sliki (kot branje - od leve proti desni, od zgoraj navzdol)

    int levi_rob = 0;
    int desni_rob = width - 1;
    int zgornji_rob = 0;
    int spodnji_rob = height - 1;
    
    for (int i = 0; i < width * height; i++) {
        int vrstica = i / width;
        int stolpec = i % width;

// TODO #2.2.3 Določimo vrednost vsem sosedom

        int zgoraj = 0;
        int spodaj = 0;
        int levo = 0;
        int desno = 0;
        int zgoraj_levo = 0;
        int zgoraj_desno = 0;
        int spodaj_levo = 0;
        int spodaj_desno = 0;

// TODO #2.2.4 Preverimo, če smo na robu slike in če nismo, določimo vrednosti pravokotnih sosedov
// TODO #2.2.5 Določimo vrednost vseh osmih sosedov
        if (vrstica != zgornji_rob) {
            zgoraj = imageIn[i - width];
        }
        if (vrstica != spodnji_rob) {
            spodaj = imageIn[i + width];
        }
        if (stolpec != levi_rob) {
            levo = imageIn[i - 1];
        }
        if (stolpec != desni_rob) {
            desno = imageIn[i + 1];
        }
        if (vrstica != zgornji_rob && stolpec != levi_rob) {
            zgoraj_levo = imageIn[i - width - 1];
        }
        if (vrstica != zgornji_rob && stolpec != desni_rob) {
            zgoraj_desno = imageIn[i - width + 1];
        }
        if (vrstica != spodnji_rob && stolpec != levi_rob) {
            spodaj_levo = imageIn[i + width - 1];
        }
        if (vrstica != spodnji_rob && stolpec != desni_rob) {
            spodaj_desno = imageIn[i + width + 1];
        }

// TODO #2.2.6 Določimo vrednost piksla glede na vrednosti sosedov
        int gx = - zgoraj_levo - 2 * levo - spodaj_levo + zgoraj_desno + 2 * desno + spodaj_desno;
        int gy = zgoraj_levo + 2 * zgoraj + zgoraj_desno - spodaj_levo - 2 * spodaj - spodaj_desno;

        double g = sqrt(gx * gx + gy * gy);

        if (g > 255) {
            g = 255;
        }

        imageOut[i] = g;

    }
}

int main(int argc, char *argv[])
{
    printf("------ Comppiled successfully ------\n");
    if (argc < 3)
    {
        printf("USAGE: sample input_image output_image\n");
        exit(EXIT_FAILURE);
    }
    
    char szImage_in_name[255];
    char szImage_out_name[255];
    char szImage_out_nameCPE[255];

    snprintf(szImage_in_name, 255, "./examples/%s", argv[1]);
    snprintf(szImage_out_name, 255, "./resultsGPE/%s", argv[2]);
    snprintf(szImage_out_nameCPE, 255, "./resultsCPE/%s", argv[2]);

    // Load image from file and allocate space for the output image
    int width, height, cpp;
    unsigned char *h_imageIn = stbi_load(szImage_in_name, &width, &height, &cpp, COLOR_CHANNELS);
    cpp = COLOR_CHANNELS;

    if (h_imageIn == NULL)
    {
        printf("Error reading loading image %s!\n", szImage_in_name);
        exit(EXIT_FAILURE);
    }
    printf("Loaded image %s of size %dx%d.\n", szImage_in_name, width, height);
    const size_t datasize = width * height * cpp * sizeof(unsigned char);
    unsigned char *h_imageOut = (unsigned char *)malloc(datasize);

    // Kot preizkus samo kopiramo vhodno sliko v izhodno
    memcpy(h_imageOut,h_imageIn,datasize);

    // Nastavimo organizacijo niti v 2D
    // dim3 blockSize(1, 1);
    dim3 blockSize(32, 32);
    // dim3 gridSize(1,1);
    dim3 gridSize((width + blockSize.x - 1) / blockSize.x, (height + blockSize.y - 1) / blockSize.y);

    unsigned char *d_imageIn;
    unsigned char *d_imageOut;

    // Rezervacija pomnilnika na napravi
    checkCudaErrors(cudaMalloc(&d_imageIn, datasize));
    checkCudaErrors(cudaMalloc(&d_imageOut, datasize));

    // Uporabimo dogodke CUDA za merjenje casa
    cudaEvent_t start, stop;
    cudaEventCreate(&start);
    cudaEventCreate(&stop);

    // Zazenemo scepec
    cudaEventRecord(start);

// TODO #1.1.1 Kopiraj prebrano črno-belo sliko v vhodno sliko na napravi
    checkCudaErrors(cudaMemcpy(d_imageIn, h_imageIn, datasize, cudaMemcpyHostToDevice));

// TODO #1.1.2 Kliči metodo za izvajanje na napravi
    printf("GPE DELA K ZMEŠAN\n");
    marginirajGPE<<<gridSize, blockSize>>>(d_imageIn, d_imageOut, width, height);
    getLastCudaError("marginirajGPE() execution failed\n");
    cudaEventRecord(stop);

    cudaEventSynchronize(stop);

    // Izpisemo cas
    float milliseconds = 0;
    cudaEventElapsedTime(&milliseconds, start, stop);

    // Zapisemo izhodno sliko v datoteko
    char szImage_out_name_temp[255];
    strncpy(szImage_out_name_temp, szImage_out_name, 255);
    char *token = strtok(szImage_out_name_temp, ".");

// TODO #2.1.1 Začni meriti čas na gostitelju
    double time_start = clock();

// TODO #2.2.1 Kliči metodo za izvajanje na gostitelju
    marginirajCPE(h_imageIn, h_imageOut, width, height);

// TODO #2.3.1 Nehaj meriti čas na gostitelju in ga izpiši
    double time_end = clock();
    double time_total = (time_end - time_start) / (CLOCKS_PER_SEC / 1000);

    printf("Kernel Execution time is: %0.3f milliseconds \n", milliseconds);
    printf("CPU Execution time is: %0.3f milliseconds \n", time_total);
    printf("Razlika: %0.3f milisekund \n", time_total - milliseconds);
    printf("GPE je porabil %0.3f% časa, ki ga je porabil CPE.\n", (milliseconds / time_total) * 100);

    char *FileType = NULL;
    while (token != NULL)
    {
        FileType = token;
        token = strtok(NULL, ".");
    }

// TODO #3.1 Glede na filetype izpišemo še rezultat gostiteljskega izvajanja v izhodno datoteko
    if (!strcmp(FileType, "png")){
        stbi_write_png(szImage_out_name, width, height, cpp, h_imageOut, width * cpp);
        stbi_write_png(szImage_out_nameCPE, width, height, cpp, h_imageOut, width * cpp);
    }else if (!strcmp(FileType, "jpg")){
        stbi_write_jpg(szImage_out_name, width, height, cpp, h_imageOut, 100);
        stbi_write_jpg(szImage_out_nameCPE, width, height, cpp, h_imageOut, 100);
    }else if (!strcmp(FileType, "bmp")){
        stbi_write_bmp(szImage_out_name, width, height, cpp, h_imageOut);
        stbi_write_bmp(szImage_out_nameCPE, width, height, cpp, h_imageOut);
    }else
        printf("Error: Unknown image format %s! Only png, bmp, or bmp supported.\n", FileType);

    // Sprostimo pomnilnik na napravi
    checkCudaErrors(cudaFree(d_imageIn));
    checkCudaErrors(cudaFree(d_imageOut));

    // Pocistimo dogodke
	cudaEventDestroy(start);
	cudaEventDestroy(stop);

    // Sprostimo pomnilnik na gostitelju
    free(h_imageIn);
    free(h_imageOut);

    return 0;
}