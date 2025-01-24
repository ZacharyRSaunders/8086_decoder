#include <stdio.h>
#include <stdlib.h>

unsigned char *read_file(const char *filename, size_t *bytesRead) {
  // Open the file
  FILE *file;
  file = fopen(filename, "rb");
  if (file == NULL) {
    perror("Error opening file");
    fclose(file);
    return NULL;
  }

  // Get the file size by seeking to the end.
  size_t fileSize;
  fseek(file, 0, SEEK_END);
  fileSize = ftell(file);
  rewind(file);

  // allocate a buffer to store the file contents.
  unsigned char *buffer;
  buffer = (unsigned char *)malloc(fileSize);
  if (buffer == NULL) {
    perror("Memory allocation failed");
    fclose(file);
    return NULL;
  }

  // Read the file contents into the buffer
  *bytesRead = fread(&buffer, 1, fileSize, file);
  if ((*bytesRead) != fileSize) {
    perror("Error: File reading Error.\n");
    free(buffer);
    fclose(file);
    return NULL;
  }
  printf("Read %zu bytes from the file\n", *bytesRead);
  fclose(file);
  return buffer;
}

// Successfully read the file and process the data

int main(int argc, char *argv[]) {
  // Read byte stream from filename passed by the arguments.
  const char *filename = argv[1];

  if (argc > 2) {
    printf("Error: Too many argumnets. Only one is allowed\n");
    return 1;
  }

  printf("Number of args: %d\n", argc);

  for (int i = 0; i < argc; i++) {
    printf("Argument %d: %s\n", i, argv[i]);
  }

  size_t bytesRead;

  unsigned char *buffer = read_file(filename, &bytesRead);
  if (buffer == NULL) {
    free(buffer);
    printf("Error: File read failed.");
  }

  // Example: Print the first few bytes
  for (size_t i = 0; i < bytesRead; i++) {
    // Get all the different parts of the first byte
    unsigned char opcode_mask = 0b11111100;
    unsigned char opcode = buffer[i] & opcode_mask;

    unsigned char wbit_mask = 0b00000001;
    unsigned char wbit = buffer[i] & wbit_mask;

    switch (opcode) {
    case 0b10001000:
      printf("mov ");

      // get the mod, reg and r/m bits
      size_t next_byte = i + 1;
      unsigned char mod_mask = 0b11000000;
      unsigned char mod = buffer[next_byte] & mod_mask;
      unsigned char reg_mask = 0b00111000;
      unsigned char reg = buffer[next_byte] & reg_mask;
      unsigned char rm_mask = 0b00000111;
      unsigned char rm = buffer[next_byte] & rm_mask;

      switch (wbit) {
      case 0b00000001:
        switch (rm) {
        case 0b00000000:
          printf("ax");
          break;
        case 0b00000001:
          printf("cx");
          break;
        case 0b00000010:
          printf("dx");
          break;
        case 0b00000011:
          printf("bx");
          break;
        case 0b00000100:
          printf("sp");
          break;
        case 0b00000101:
          printf("bp");
          break;
        case 0b00000110:
          printf("si");
          break;
        case 0b00000111:
          printf("di");
          break;
        default:
          printf("Error\n");
          free(buffer);
          return 1;
        }
        printf(", ");
        switch (reg) {
        case 0b00000000:
          printf("ax");
          break;
        case 0b00001000:
          printf("cx");
          break;
        case 0b00010000:
          printf("dx");
          break;
        case 0b00011000:
          printf("bx");
          break;
        case 0b00100000:
          printf("sp");
          break;
        case 0b00101000:
          printf("bp");
          break;
        case 0b00110000:
          printf("si");
          break;
        case 0b00111000:
          printf("di");
          break;
        default:
          printf("Error\n");
          free(buffer);
          return 1;
        }
        printf("\n");
        break;
      default:
        switch (rm) {
        case 0b00000000:
          printf("al");
          break;
        case 0b00000001:
          printf("cl");
          break;
        case 0b00000010:
          printf("dl");
          break;
        case 0b00000011:
          printf("bl");
          break;
        case 0b00000100:
          printf("ah");
          break;
        case 0b00000101:
          printf("ch");
          break;
        case 0b00000110:
          printf("dh");
          break;
        case 0b00000111:
          printf("bh");
          break;
        default:
          printf("Error\n");
          free(buffer);
          return 1;
        }
        printf(", ");
        switch (reg) {
        case 0b00000000:
          printf("al");
          break;
        case 0b00001000:
          printf("cl");
          break;
        case 0b00010000:
          printf("dl");
          break;
        case 0b00011000:
          printf("bl");
          break;
        case 0b00100000:
          printf("ah");
          break;
        case 0b00101000:
          printf("ch");
          break;
        case 0b00110000:
          printf("dh");
          break;
        case 0b00111000:
          printf("bh");
          break;
        default:
          printf("Error\n");
          free(buffer);
          return 1;
        }
        printf("\n");
      }
      i++;
      break;
    default:
      printf("Error: Opcode not implemented.");
      free(buffer);
    }
  }
  printf("\n");

  // Loop through the bytes and decode each one.
  // for (size_t i = 0; i < bytesRead; i++) {
  //   if printf ("%b ", buffer[i])
  //     ;
  // }

  // Create some sort of switch for each combination of bits.

  // clean unsigned
  free(buffer);
  return 0;
}
