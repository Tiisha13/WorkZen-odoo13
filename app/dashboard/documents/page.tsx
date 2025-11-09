"use client";

import { useRequireAuth, usePageTitle } from "@/lib/hooks";
import { useState, useEffect } from "react";
import { apiService } from "@/lib/api-service";
import { API_ENDPOINTS } from "@/lib/config";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Textarea } from "@/components/ui/textarea";
import {
  IconPlus,
  IconDownload,
  IconEye,
  IconSearch,
} from "@tabler/icons-react";
import { toast } from "sonner";
import type { Document } from "@/lib/types";

export default function DocumentsPage() {
  usePageTitle("Document Management | WorkZen");
  useRequireAuth();

  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [viewDialogOpen, setViewDialogOpen] = useState(false);
  const [selectedDocument, setSelectedDocument] = useState<Document | null>(
    null
  );
  const [documentUrl, setDocumentUrl] = useState<string>("");

  const [formData, setFormData] = useState({
    description: "",
    category: "other",
    employee_id: "",
  });

  useEffect(() => {
    fetchDocuments();
  }, []);

  const fetchDocuments = async () => {
    try {
      setLoading(true);
      const response = await apiService.get<{
        success: boolean;
        data: Document[];
      }>(API_ENDPOINTS.DOCUMENTS);
      setDocuments(response.data || []);
    } catch (error) {
      toast.error("Failed to fetch documents");
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedFile) {
      toast.error("Please select a file");
      return;
    }

    try {
      const formDataToSend = new FormData();
      formDataToSend.append("file", selectedFile);
      formDataToSend.append("category", formData.category);
      formDataToSend.append("description", formData.description);
      if (formData.employee_id) {
        formDataToSend.append("employee_id", formData.employee_id);
      }

      const response = await fetch(API_ENDPOINTS.DOCUMENTS, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: formDataToSend,
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "Failed to upload document");
      }

      toast.success("Document uploaded successfully");
      setIsDialogOpen(false);
      resetForm();
      fetchDocuments();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to upload document"
      );
    }
  };

  const isViewable = (fileType: string) => {
    return (
      fileType?.startsWith("image/") ||
      fileType?.startsWith("video/") ||
      fileType?.startsWith("audio/") ||
      fileType === "application/pdf"
    );
  };

  const handleView = async (doc: Document) => {
    try {
      setSelectedDocument(doc);
      const token = localStorage.getItem("token");

      // Fetch the document with auth header
      const response = await fetch(
        `${API_ENDPOINTS.DOCUMENTS}/${doc.id}/view`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        const errorMessage =
          errorData.message || `Failed to load document (${response.status})`;
        throw new Error(errorMessage);
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      setDocumentUrl(url);
      setViewDialogOpen(true);
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to view document";
      toast.error(message);
      console.error("View document error:", error);
    }
  };

  const handleDownload = async (doc: Document) => {
    try {
      toast.info("Downloading document...");
      const token = localStorage.getItem("token");

      const response = await fetch(
        `${API_ENDPOINTS.DOCUMENTS}/${doc.id}/download`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to download document");
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = doc.file_name || "document";
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      toast.success("Document downloaded successfully");
    } catch (error) {
      toast.error("Failed to download document");
      console.error(error);
    }
  };

  const handleCloseViewDialog = () => {
    setViewDialogOpen(false);
    if (documentUrl) {
      URL.revokeObjectURL(documentUrl);
      setDocumentUrl("");
    }
    setSelectedDocument(null);
  };

  const resetForm = () => {
    setFormData({
      description: "",
      category: "other",
      employee_id: "",
    });
    setSelectedFile(null);
  };

  const filteredDocuments = documents.filter(
    (doc) =>
      doc.file_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      doc.category.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + " B";
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + " KB";
    return (bytes / (1024 * 1024)).toFixed(2) + " MB";
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  const getCategoryBadge = (category: string) => {
    const colors: Record<string, string> = {
      general: "bg-gray-100 text-gray-800",
      contract: "bg-blue-100 text-blue-800",
      policy: "bg-purple-100 text-purple-800",
      payslip: "bg-green-100 text-green-800",
      certificate: "bg-yellow-100 text-yellow-800",
    };
    return (
      <Badge className={colors[category] || colors.general}>
        {category.toUpperCase()}
      </Badge>
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Document Management</h1>
          <p className="text-muted-foreground">Upload and manage documents</p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={resetForm}>
              <IconPlus className="w-4 h-4 mr-2" />
              Upload Document
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Upload New Document</DialogTitle>
              <DialogDescription>
                Upload a new document to the system
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="category">Category *</Label>
                <select
                  id="category"
                  value={formData.category}
                  onChange={(e) =>
                    setFormData({ ...formData, category: e.target.value })
                  }
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                  required
                >
                  <option value="resume">Resume</option>
                  <option value="id_proof">ID Proof</option>
                  <option value="payslip">Payslip</option>
                  <option value="policy">Policy</option>
                  <option value="report">Report</option>
                  <option value="other">Other</option>
                </select>
                <p className="text-xs text-muted-foreground">
                  Payslip documents are automatically marked as private
                </p>
              </div>
              <div className="space-y-2">
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) =>
                    setFormData({ ...formData, description: e.target.value })
                  }
                  rows={3}
                  placeholder="Optional notes about this document"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="file">File *</Label>
                <Input
                  id="file"
                  type="file"
                  onChange={handleFileChange}
                  required
                />
                {selectedFile && (
                  <p className="text-sm text-muted-foreground">
                    Selected: {selectedFile.name} (
                    {formatFileSize(selectedFile.size)})
                  </p>
                )}
              </div>
              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsDialogOpen(false)}
                >
                  Cancel
                </Button>
                <Button type="submit">Upload Document</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        {/* View Document Dialog */}
        <Dialog open={viewDialogOpen} onOpenChange={handleCloseViewDialog}>
          <DialogContent className="max-w-4xl max-h-[90vh] overflow-auto">
            <DialogHeader>
              <DialogTitle>
                {selectedDocument?.file_name || "Document"}
              </DialogTitle>
              <DialogDescription>
                Category: {selectedDocument?.category} | Size:{" "}
                {selectedDocument?.size
                  ? formatFileSize(selectedDocument.size)
                  : "N/A"}
              </DialogDescription>
            </DialogHeader>
            <div className="mt-4">
              {selectedDocument && documentUrl && (
                <>
                  {selectedDocument.file_type?.startsWith("image/") && (
                    // eslint-disable-next-line @next/next/no-img-element
                    <img
                      src={documentUrl}
                      alt={selectedDocument.file_name}
                      className="w-full h-auto rounded-lg"
                    />
                  )}
                  {selectedDocument.file_type?.startsWith("video/") && (
                    <video
                      src={documentUrl}
                      controls
                      className="w-full rounded-lg"
                    >
                      Your browser does not support the video tag.
                    </video>
                  )}
                  {selectedDocument.file_type?.startsWith("audio/") && (
                    <audio src={documentUrl} controls className="w-full">
                      Your browser does not support the audio tag.
                    </audio>
                  )}
                  {selectedDocument.file_type === "application/pdf" && (
                    <iframe
                      src={documentUrl}
                      className="w-full h-[70vh] rounded-lg"
                      title={selectedDocument.file_name}
                    />
                  )}
                  {!isViewable(selectedDocument.file_type || "") && (
                    <div className="text-center py-8">
                      <p className="text-muted-foreground mb-4">
                        This file type cannot be previewed. Please download to
                        view.
                      </p>
                      <Button onClick={() => handleDownload(selectedDocument)}>
                        <IconDownload className="w-4 h-4 mr-2" />
                        Download File
                      </Button>
                    </div>
                  )}
                </>
              )}
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => handleDownload(selectedDocument!)}
              >
                <IconDownload className="w-4 h-4 mr-2" />
                Download
              </Button>
              <Button variant="outline" onClick={handleCloseViewDialog}>
                Close
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="flex items-center space-x-2">
        <div className="relative flex-1 max-w-sm">
          <IconSearch className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input
            placeholder="Search documents..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10"
          />
        </div>
      </div>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Title</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>File Name</TableHead>
              <TableHead>File Size</TableHead>
              <TableHead>Uploaded Date</TableHead>
              <TableHead>Privacy</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center py-8">
                  Loading...
                </TableCell>
              </TableRow>
            ) : filteredDocuments.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={7}
                  className="text-center py-8 text-muted-foreground"
                >
                  No documents found
                </TableCell>
              </TableRow>
            ) : (
              filteredDocuments.map((doc) => (
                <TableRow key={doc.id}>
                  <TableCell className="font-medium">
                    {doc.description || doc.file_name}
                  </TableCell>
                  <TableCell>{getCategoryBadge(doc.category)}</TableCell>
                  <TableCell>{doc.file_name}</TableCell>
                  <TableCell>{formatFileSize(doc.size)}</TableCell>
                  <TableCell>
                    {doc.created_at ? formatDate(doc.created_at) : "-"}
                  </TableCell>
                  <TableCell>
                    <Badge variant={doc.is_private ? "secondary" : "outline"}>
                      {doc.is_private ? "Private" : "Public"}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right space-x-2">
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => handleView(doc)}
                    >
                      <IconEye className="w-4 h-4" />
                    </Button>
                    <Button
                      size="sm"
                      variant="ghost"
                      onClick={() => handleDownload(doc)}
                    >
                      <IconDownload className="w-4 h-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
