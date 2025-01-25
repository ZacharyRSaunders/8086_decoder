#include <stdio.h>
#include <stdlib.h>

// Opens a file and returns a buffer
unsigned char *read_file(const char *filename, size_t *bytesRead) {
  // Open the file
  FILE *file;
  printf("INFO: Opening %s\n", filename);
  file = fopen(filename, "rb");
  if (file == NULL) {
    perror("Error opening file");
    fclose(file);
    return NULL;
  }
  printf("File Opened Successfully\n");

  // Get the file size by seeking to the end.
  size_t fileSize;
  fseek(file, 0, SEEK_END);
  fileSize = ftell(file);
  rewind(file);

  // allocate a buffer to store the file contents.
  unsigned char *buffer = (unsigned char *)malloc(fileSize);
  printf("INFO: Assigned buffer\n");
  if (buffer == NULL) {
    perror("Memory allocation failed\n");
    fclose(file);
    return NULL;
  }

  // Read the file contents into the buffer
  size_t bytesRead1 = fread(buffer, 1, fileSize, file);
  if (bytesRead1 != fileSize) {
    perror("Error: File reading Error from me.\n");
    free(buffer);
    fclose(file);
    return NULL;
  }
  printf("Read %zu bytes from the file\n", bytesRead1);
  *bytesRead = bytesRead1;
  fclose(file);
  return buffer;
}

// Map bits to registers.  Checks the last 3 bits as per the 8067 instruction
// manual.
char *map_register(unsigned char reg, unsigned char wbit) {
  if (wbit) {
    switch (reg) {
    case 0b00000000:
      return "ax";
      break;
    case 0b00000001:
      return "cx";
      break;
    case 0b00000010:
      return "dx";
      break;
    case 0b00000011:
      return "bx";
      break;
    case 0b00000100:
      return "sp";
      break;
    case 0b00000101:
      return "bp";
      break;
    case 0b00000110:
      return "si";
      break;
    case 0b00000111:
      return "di";
      break;
    default:
      printf("Error\n");
      return NULL;
    }
  } else {
    switch (reg) {
    case 0b00000000:
      return "al";
      break;
    case 0b00000001:
      return "cl";
      break;
    case 0b00000010:
      return "dl";
      break;
    case 0b00000011:
      return "bl";
      break;
    case 0b00000100:
      return "ah";
      break;
    case 0b00000101:
      return "ch";
      break;
    case 0b00000110:
      return "dh";
      break;
    case 0b00000111:
      return "bh";
      break;
    default:
      printf("Error\n");
      return NULL;
    }
  }
}

int main(int argc, char *argv[]) {
  // Read byte stream from filename passed by the arguments.
  const char *filename = argv[1];

  if (argc > 2) {
    printf("Error: Too many argumnets. Only one is allowed\n");
    return 1;
  }

  size_t bytesRead = 0;

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
      printf("  ");
      // get the mod, reg and r/m bits
      size_t next_byte = i + 1;
      unsigned char mod_mask = 0b11000000;
      unsigned char mod = buffer[next_byte] & mod_mask;
      unsigned char reg_mask = 0b00111000;
      unsigned char reg = buffer[next_byte] & reg_mask;
      unsigned char reg_shifted = reg >> 3;
      unsigned char rm_mask = 0b00000111;
      unsigned char rm = buffer[next_byte] & rm_mask;

      printf("mov %s, %s\n", map_register(reg_shifted, wbit),
             map_register(rm, wbit));
      i++;
      break;
    default:
      printf("Error: Opcode not implemented.");
      free(buffer);
    }
  }
  printf("\n");

  // clean up
  free(buffer);
  return 0;
}
