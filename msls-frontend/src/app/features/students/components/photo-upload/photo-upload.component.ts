/**
 * Photo Upload Component
 *
 * Handles photo selection, validation, compression, cropping, and preview for student profiles.
 * Features:
 * - Auto-compression for images larger than maxSizeMB
 * - Built-in image cropper with zoom and rotate
 * - Drag & drop support
 */

import {
  Component,
  ChangeDetectionStrategy,
  input,
  output,
  signal,
  computed,
  ElementRef,
  viewChild,
  inject,
} from '@angular/core';
import { CommonModule } from '@angular/common';

import { ToastService } from '../../../../shared/services';

/** Allowed image MIME types */
const ALLOWED_TYPES = ['image/jpeg', 'image/png', 'image/jpg'];

/** Maximum dimension for compressed images */
const MAX_DIMENSION = 800;

/** Target quality for JPEG compression */
const COMPRESSION_QUALITY = 0.8;

interface CropArea {
  x: number;
  y: number;
  width: number;
  height: number;
}

@Component({
  selector: 'msls-photo-upload',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './photo-upload.component.html',
  styleUrls: ['./photo-upload.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PhotoUploadComponent {
  private toast = inject(ToastService);

  // =========================================================================
  // Inputs
  // =========================================================================

  /** Current preview URL (for edit mode) */
  readonly previewUrl = input<string | null>(null);

  /** Maximum file size in MB */
  readonly maxSizeMB = input<number>(2);

  /** Custom accepted file types */
  readonly acceptTypes = ALLOWED_TYPES.join(',');

  // =========================================================================
  // Outputs
  // =========================================================================

  /** Emitted when a valid photo is selected */
  readonly photoSelected = output<File>();

  /** Emitted when photo is removed */
  readonly photoRemoved = output<void>();

  // =========================================================================
  // State
  // =========================================================================

  readonly isDragging = signal<boolean>(false);
  readonly fileInput = viewChild<ElementRef<HTMLInputElement>>('fileInput');
  readonly cropperCanvas = viewChild<ElementRef<HTMLCanvasElement>>('cropperCanvas');
  readonly previewCanvas = viewChild<ElementRef<HTMLCanvasElement>>('previewCanvas');

  // Cropper state
  readonly showCropper = signal<boolean>(false);
  readonly cropperImage = signal<HTMLImageElement | null>(null);
  readonly originalFile = signal<File | null>(null);
  readonly isProcessing = signal<boolean>(false);
  readonly rotation = signal<number>(0);
  readonly zoom = signal<number>(1);
  readonly cropArea = signal<CropArea>({ x: 0, y: 0, width: 200, height: 200 });

  // Drag state for crop area
  private isDraggingCrop = false;
  private isResizingCrop = false;
  private resizeCorner: 'tl' | 'tr' | 'bl' | 'br' | null = null;
  private dragStartX = 0;
  private dragStartY = 0;
  private cropStartX = 0;
  private cropStartY = 0;
  private cropStartWidth = 0;
  private cropStartHeight = 0;

  // =========================================================================
  // Computed
  // =========================================================================

  readonly maxSizeBytes = computed(() => this.maxSizeMB() * 1024 * 1024);

  // =========================================================================
  // Methods
  // =========================================================================

  triggerFileInput(): void {
    this.fileInput()?.nativeElement.click();
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];

    if (file) {
      this.processFile(file);
    }

    // Reset input to allow selecting same file again
    input.value = '';
  }

  onRemove(): void {
    this.photoRemoved.emit();
  }

  onEditPhoto(): void {
    // If we have a preview URL, load it for editing
    const url = this.previewUrl();
    if (url) {
      this.loadImageForCropping(url);
    }
  }

  // =========================================================================
  // Drag & Drop
  // =========================================================================

  onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging.set(true);
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging.set(false);
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging.set(false);

    const file = event.dataTransfer?.files?.[0];
    if (file) {
      this.processFile(file);
    }
  }

  // =========================================================================
  // File Processing
  // =========================================================================

  private processFile(file: File): void {
    // Validate file type
    if (!ALLOWED_TYPES.includes(file.type)) {
      this.toast.error('Invalid file type. Please upload a JPG or PNG image.');
      return;
    }

    // Store original file and open cropper
    this.originalFile.set(file);
    const url = URL.createObjectURL(file);
    this.loadImageForCropping(url);
  }

  private loadImageForCropping(url: string): void {
    const img = new Image();
    img.onload = () => {
      this.cropperImage.set(img);
      this.rotation.set(0);
      this.zoom.set(1);
      this.showCropper.set(true);

      // Initialize crop area centered
      setTimeout(() => this.initializeCropArea(), 50);
    };
    img.onerror = () => {
      this.toast.error('Failed to load image. Please try another file.');
    };
    img.src = url;
  }

  private initializeCropArea(): void {
    const canvas = this.cropperCanvas()?.nativeElement;
    const img = this.cropperImage();
    if (!canvas || !img) return;

    // Set canvas size to fit the container
    const container = canvas.parentElement;
    if (container) {
      canvas.width = container.clientWidth || 400;
      canvas.height = container.clientHeight || 400;
    }

    // Calculate initial crop area (centered square)
    const size = Math.min(canvas.width, canvas.height) * 0.7;
    this.cropArea.set({
      x: (canvas.width - size) / 2,
      y: (canvas.height - size) / 2,
      width: size,
      height: size,
    });

    this.drawCropper();
  }

  // =========================================================================
  // Cropper Controls
  // =========================================================================

  onZoomIn(): void {
    this.zoom.update((z) => Math.min(z + 0.1, 3));
    this.drawCropper();
  }

  onZoomOut(): void {
    this.zoom.update((z) => Math.max(z - 0.1, 0.5));
    this.drawCropper();
  }

  onRotateLeft(): void {
    this.rotation.update((r) => r - 90);
    this.drawCropper();
  }

  onRotateRight(): void {
    this.rotation.update((r) => r + 90);
    this.drawCropper();
  }

  onResetCrop(): void {
    this.rotation.set(0);
    this.zoom.set(1);
    this.initializeCropArea();
  }

  onCancelCrop(): void {
    this.showCropper.set(false);
    this.cropperImage.set(null);
    this.originalFile.set(null);
  }

  async onApplyCrop(): Promise<void> {
    this.isProcessing.set(true);

    try {
      const croppedBlob = await this.getCroppedImage();
      if (croppedBlob) {
        // Compress if needed
        const finalBlob = await this.compressIfNeeded(croppedBlob);
        const file = new File([finalBlob], this.originalFile()?.name || 'photo.jpg', {
          type: 'image/jpeg',
        });

        this.photoSelected.emit(file);
        this.showCropper.set(false);
        this.cropperImage.set(null);
        this.originalFile.set(null);

        this.toast.success('Photo updated successfully!');
      }
    } catch {
      this.toast.error('Failed to process image. Please try again.');
    } finally {
      this.isProcessing.set(false);
    }
  }

  // =========================================================================
  // Canvas Drawing
  // =========================================================================

  drawCropper(): void {
    const canvas = this.cropperCanvas()?.nativeElement;
    const img = this.cropperImage();
    if (!canvas || !img) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const zoom = this.zoom();
    const rotation = this.rotation();
    const crop = this.cropArea();

    // Clear canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Save context
    ctx.save();

    // Move to center
    ctx.translate(canvas.width / 2, canvas.height / 2);

    // Apply rotation
    ctx.rotate((rotation * Math.PI) / 180);

    // Apply zoom
    ctx.scale(zoom, zoom);

    // Calculate scaled dimensions to fit canvas
    const scale = Math.min(
      (canvas.width * 0.9) / img.width,
      (canvas.height * 0.9) / img.height
    );
    const scaledWidth = img.width * scale;
    const scaledHeight = img.height * scale;

    // Draw image centered
    ctx.drawImage(img, -scaledWidth / 2, -scaledHeight / 2, scaledWidth, scaledHeight);

    // Restore context
    ctx.restore();

    // Draw dark overlay outside crop area
    ctx.fillStyle = 'rgba(0, 0, 0, 0.5)';
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Clear the crop area (make it visible)
    ctx.clearRect(crop.x, crop.y, crop.width, crop.height);

    // Redraw image in crop area only
    ctx.save();
    ctx.beginPath();
    ctx.rect(crop.x, crop.y, crop.width, crop.height);
    ctx.clip();

    ctx.translate(canvas.width / 2, canvas.height / 2);
    ctx.rotate((rotation * Math.PI) / 180);
    ctx.scale(zoom, zoom);
    ctx.drawImage(img, -scaledWidth / 2, -scaledHeight / 2, scaledWidth, scaledHeight);
    ctx.restore();

    // Draw crop area border
    ctx.strokeStyle = '#ffffff';
    ctx.lineWidth = 2;
    ctx.strokeRect(crop.x, crop.y, crop.width, crop.height);

    // Draw corner handles
    this.drawCornerHandles(ctx, crop);

    // Draw grid lines
    this.drawGridLines(ctx, crop);

    // Update preview
    this.updatePreview();
  }

  private drawCornerHandles(ctx: CanvasRenderingContext2D, crop: CropArea): void {
    const handleSize = 10;
    ctx.fillStyle = '#ffffff';

    // Corner positions
    const corners = [
      { x: crop.x, y: crop.y }, // top-left
      { x: crop.x + crop.width, y: crop.y }, // top-right
      { x: crop.x, y: crop.y + crop.height }, // bottom-left
      { x: crop.x + crop.width, y: crop.y + crop.height }, // bottom-right
    ];

    corners.forEach((corner) => {
      ctx.fillRect(
        corner.x - handleSize / 2,
        corner.y - handleSize / 2,
        handleSize,
        handleSize
      );
    });
  }

  private drawGridLines(ctx: CanvasRenderingContext2D, crop: CropArea): void {
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.3)';
    ctx.lineWidth = 1;

    // Vertical lines (rule of thirds)
    for (let i = 1; i < 3; i++) {
      const x = crop.x + (crop.width * i) / 3;
      ctx.beginPath();
      ctx.moveTo(x, crop.y);
      ctx.lineTo(x, crop.y + crop.height);
      ctx.stroke();
    }

    // Horizontal lines
    for (let i = 1; i < 3; i++) {
      const y = crop.y + (crop.height * i) / 3;
      ctx.beginPath();
      ctx.moveTo(crop.x, y);
      ctx.lineTo(crop.x + crop.width, y);
      ctx.stroke();
    }
  }

  private updatePreview(): void {
    const previewCanvas = this.previewCanvas()?.nativeElement;
    const mainCanvas = this.cropperCanvas()?.nativeElement;
    if (!previewCanvas || !mainCanvas) return;

    const ctx = previewCanvas.getContext('2d');
    if (!ctx) return;

    const crop = this.cropArea();

    // Set preview size
    previewCanvas.width = 120;
    previewCanvas.height = 120;

    // Draw circular clip
    ctx.beginPath();
    ctx.arc(60, 60, 60, 0, Math.PI * 2);
    ctx.closePath();
    ctx.clip();

    // Draw the cropped area
    ctx.drawImage(
      mainCanvas,
      crop.x,
      crop.y,
      crop.width,
      crop.height,
      0,
      0,
      120,
      120
    );
  }

  // =========================================================================
  // Crop Area Interaction
  // =========================================================================

  onCropMouseDown(event: MouseEvent): void {
    const canvas = this.cropperCanvas()?.nativeElement;
    if (!canvas) return;

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;
    const crop = this.cropArea();

    // Check if clicking on a corner (resize)
    const handleSize = 18;
    const corners: Array<{ x: number; y: number; corner: 'tl' | 'tr' | 'bl' | 'br' }> = [
      { x: crop.x, y: crop.y, corner: 'tl' }, // top-left
      { x: crop.x + crop.width, y: crop.y, corner: 'tr' }, // top-right
      { x: crop.x, y: crop.y + crop.height, corner: 'bl' }, // bottom-left
      { x: crop.x + crop.width, y: crop.y + crop.height, corner: 'br' }, // bottom-right
    ];

    for (const c of corners) {
      if (
        Math.abs(x - c.x) < handleSize &&
        Math.abs(y - c.y) < handleSize
      ) {
        this.isResizingCrop = true;
        this.resizeCorner = c.corner;
        this.dragStartX = x;
        this.dragStartY = y;
        this.cropStartX = crop.x;
        this.cropStartY = crop.y;
        this.cropStartWidth = crop.width;
        this.cropStartHeight = crop.height;
        return;
      }
    }

    // Check if clicking inside crop area (move)
    if (
      x >= crop.x &&
      x <= crop.x + crop.width &&
      y >= crop.y &&
      y <= crop.y + crop.height
    ) {
      this.isDraggingCrop = true;
      this.dragStartX = x;
      this.dragStartY = y;
      this.cropStartX = crop.x;
      this.cropStartY = crop.y;
    }
  }

  onCropMouseMove(event: MouseEvent): void {
    const canvas = this.cropperCanvas()?.nativeElement;
    if (!canvas) return;

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    if (this.isDraggingCrop) {
      const dx = x - this.dragStartX;
      const dy = y - this.dragStartY;

      const crop = this.cropArea();
      let newX = this.cropStartX + dx;
      let newY = this.cropStartY + dy;

      // Constrain to canvas bounds
      newX = Math.max(0, Math.min(newX, canvas.width - crop.width));
      newY = Math.max(0, Math.min(newY, canvas.height - crop.height));

      this.cropArea.set({ ...crop, x: newX, y: newY });
      this.drawCropper();
    } else if (this.isResizingCrop && this.resizeCorner) {
      const dx = x - this.dragStartX;
      const dy = y - this.dragStartY;

      let newX = this.cropStartX;
      let newY = this.cropStartY;
      let newSize = this.cropStartWidth;

      // Calculate new size based on which corner is being dragged
      switch (this.resizeCorner) {
        case 'br': // bottom-right: grow/shrink from top-left anchor
          newSize = this.cropStartWidth + (dx + dy) / 2;
          break;
        case 'bl': // bottom-left: grow/shrink from top-right anchor
          newSize = this.cropStartWidth + (-dx + dy) / 2;
          newX = this.cropStartX + this.cropStartWidth - newSize;
          break;
        case 'tr': // top-right: grow/shrink from bottom-left anchor
          newSize = this.cropStartWidth + (dx - dy) / 2;
          newY = this.cropStartY + this.cropStartHeight - newSize;
          break;
        case 'tl': // top-left: grow/shrink from bottom-right anchor
          newSize = this.cropStartWidth + (-dx - dy) / 2;
          newX = this.cropStartX + this.cropStartWidth - newSize;
          newY = this.cropStartY + this.cropStartHeight - newSize;
          break;
      }

      // Enforce minimum size
      const minSize = 50;
      if (newSize < minSize) {
        const diff = minSize - newSize;
        newSize = minSize;
        // Adjust position to keep the anchor corner fixed
        if (this.resizeCorner === 'tl' || this.resizeCorner === 'bl') {
          newX -= diff;
        }
        if (this.resizeCorner === 'tl' || this.resizeCorner === 'tr') {
          newY -= diff;
        }
      }

      // Constrain to canvas bounds
      newX = Math.max(0, newX);
      newY = Math.max(0, newY);
      newSize = Math.min(newSize, canvas.width - newX, canvas.height - newY);

      this.cropArea.set({ x: newX, y: newY, width: newSize, height: newSize });
      this.drawCropper();
    }
  }

  onCropMouseUp(): void {
    this.isDraggingCrop = false;
    this.isResizingCrop = false;
    this.resizeCorner = null;
  }

  // =========================================================================
  // Image Processing
  // =========================================================================

  private async getCroppedImage(): Promise<Blob | null> {
    const canvas = this.cropperCanvas()?.nativeElement;
    const img = this.cropperImage();
    if (!canvas || !img) return null;

    const crop = this.cropArea();
    const zoom = this.zoom();
    const rotation = this.rotation();

    // Create output canvas
    const outputCanvas = document.createElement('canvas');
    const outputSize = Math.min(MAX_DIMENSION, crop.width);
    outputCanvas.width = outputSize;
    outputCanvas.height = outputSize;

    const ctx = outputCanvas.getContext('2d');
    if (!ctx) return null;

    // Calculate the scale used in the main canvas
    const scale = Math.min(
      (canvas.width * 0.9) / img.width,
      (canvas.height * 0.9) / img.height
    );

    // Calculate source coordinates in the original image
    const centerX = canvas.width / 2;
    const centerY = canvas.height / 2;

    // Create a temporary canvas for rotation
    const tempCanvas = document.createElement('canvas');
    const tempCtx = tempCanvas.getContext('2d');
    if (!tempCtx) return null;

    // Set temp canvas size
    const maxDim = Math.max(img.width, img.height) * Math.sqrt(2);
    tempCanvas.width = maxDim;
    tempCanvas.height = maxDim;

    // Draw rotated image to temp canvas
    tempCtx.translate(maxDim / 2, maxDim / 2);
    tempCtx.rotate((rotation * Math.PI) / 180);
    tempCtx.scale(zoom, zoom);
    tempCtx.drawImage(img, -img.width / 2, -img.height / 2);

    // Calculate source position relative to the crop area
    const sourceX = ((crop.x - centerX) / scale / zoom) + maxDim / 2;
    const sourceY = ((crop.y - centerY) / scale / zoom) + maxDim / 2;
    const sourceSize = crop.width / scale / zoom;

    // Draw cropped area to output
    ctx.drawImage(
      tempCanvas,
      sourceX,
      sourceY,
      sourceSize,
      sourceSize,
      0,
      0,
      outputSize,
      outputSize
    );

    return new Promise((resolve) => {
      outputCanvas.toBlob(
        (blob) => resolve(blob),
        'image/jpeg',
        COMPRESSION_QUALITY
      );
    });
  }

  private async compressIfNeeded(blob: Blob): Promise<Blob> {
    if (blob.size <= this.maxSizeBytes()) {
      return blob;
    }

    // Need to compress
    const img = await this.blobToImage(blob);
    let quality = COMPRESSION_QUALITY;
    let result = blob;

    // Iteratively reduce quality until size is acceptable
    while (result.size > this.maxSizeBytes() && quality > 0.1) {
      quality -= 0.1;
      result = await this.compressImage(img, quality);
    }

    // If still too large, reduce dimensions
    if (result.size > this.maxSizeBytes()) {
      result = await this.reduceImageDimensions(img);
    }

    const originalSizeKB = Math.round(blob.size / 1024);
    const newSizeKB = Math.round(result.size / 1024);
    if (newSizeKB < originalSizeKB) {
      this.toast.info(`Image compressed from ${originalSizeKB}KB to ${newSizeKB}KB`);
    }

    return result;
  }

  private blobToImage(blob: Blob): Promise<HTMLImageElement> {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.onload = () => resolve(img);
      img.onerror = reject;
      img.src = URL.createObjectURL(blob);
    });
  }

  private compressImage(img: HTMLImageElement, quality: number): Promise<Blob> {
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;

    const ctx = canvas.getContext('2d');
    if (!ctx) {
      return Promise.resolve(new Blob());
    }

    ctx.drawImage(img, 0, 0);

    return new Promise((resolve) => {
      canvas.toBlob(
        (blob) => resolve(blob || new Blob()),
        'image/jpeg',
        quality
      );
    });
  }

  private async reduceImageDimensions(img: HTMLImageElement): Promise<Blob> {
    const canvas = document.createElement('canvas');
    const ratio = Math.min(MAX_DIMENSION / img.width, MAX_DIMENSION / img.height);

    canvas.width = img.width * ratio;
    canvas.height = img.height * ratio;

    const ctx = canvas.getContext('2d');
    if (!ctx) {
      return new Blob();
    }

    ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

    return new Promise((resolve) => {
      canvas.toBlob(
        (blob) => resolve(blob || new Blob()),
        'image/jpeg',
        COMPRESSION_QUALITY
      );
    });
  }
}
