import { apiClient } from "./client";
import type { ApiResponse } from "@/types/api";
import type { Favorite } from "@/types/favorite";

export async function getMyFavorites(): Promise<ApiResponse<Favorite[]>> {
  const { data } = await apiClient.get<ApiResponse<Favorite[]>>("/favorites");
  return data;
}

export async function addFavorite(
  hotelId: number | string
): Promise<ApiResponse<Favorite>> {
  const { data } = await apiClient.post<ApiResponse<Favorite>>(
    `/hotels/${hotelId}/favorite`
  );
  return data;
}

export async function removeFavorite(
  hotelId: number | string
): Promise<ApiResponse<undefined>> {
  const { data } = await apiClient.delete<ApiResponse<undefined>>(
    `/hotels/${hotelId}/favorite`
  );
  return data;
}
